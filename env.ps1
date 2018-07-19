$scriptPath = split-path -parent $MyInvocation.MyCommand.Definition
$env:Path += ";" + $scriptPath + "\\bin"
