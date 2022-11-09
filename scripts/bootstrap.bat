@echo off

cmd /c "cd ../web & yarn install"

IF NOT EXIST "../web/.env" (
    (
        echo SKIP_PREFLIGHT_CHECK=true
        echo BROWSER=none
    ) > ../web/.env
    @echo Env file not found. One was created for you.
)

@echo Bootstrapping complete.