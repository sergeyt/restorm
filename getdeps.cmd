for /F "tokens=*" %%i in (deps) do go get -u %%i
