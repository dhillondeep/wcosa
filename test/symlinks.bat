@echo off
set DIR=%~dp0
cd %DIR%
goto check_permissions

:check_permissions
net session >nul 2>&1
if NOT %errorLevel% == 0 (
    echo Admin rights are needed to create symlinks
    exit 1
)

mkdir %DIR%\project-pkg\pkg-list\vendor 2>NUL
rmdir %DIR%\project-pkg\pkg-list\vendor\pkg-malloc 2>NUL
mklink /D %DIR%\project-pkg\pkg-list\vendor\pkg-malloc %DIR%\project-pkg\pkg-malloc

mkdir %DIR%\project-app\app-avr\vendor 2>NUL
rmdir %DIR%\project-app\app-avr\vendor\pkg-uart 2>NUL
rmdir %DIR%\project-app\app-avr\vendor\pkg-list 2>NUL
rmdir %DIR%\project-app\app-avr\vendor\pkg-malloc 2>NUL
mklink /D %DIR%\project-app\app-avr\vendor\pkg-uart %DIR%\project-pkg\pkg-uart
mklink /D %DIR%\project-app\app-avr\vendor\pkg-list %DIR%\project-pkg\pkg-list
mklink /D %DIR%\project-app\app-avr\vendor\pkg-malloc %DIR%\project-pkg\pkg-malloc

mkdir %DIR%\project-pkg\pkg-trace\vendor 2>NUL
rmdir %DIR%\project-pkg\pkg-trace\vendor\pkg-uart 2>NUL
mklink /D %DIR%\project-pkg\pkg-trace\vendor\pkg-uart %DIR%\project-pkg\pkg-uart

mkdir %DIR%\project-app\app-alloc\vendor 2>NUL
mkdir %DIR%\project-app\app-alloc\vendor\alloc-one\vendor 2>NUL
mkdir %DIR%\project-app\app-alloc\vendor\alloc-two\vendor 2>NUL
rmdir %DIR%\project-app\app-alloc\vendor\pkg-malloc 2>NUL
rmdir %DIR%\project-app\app-alloc\vendor\alloc-one\vendor\pkg-malloc 2>NUL
rmdir %DIR%\project-app\app-alloc\vendor\alloc-two\vendor\pkg-malloc 2>NUL
mklink /D %DIR%\project-app\app-alloc\vendor\pkg-malloc %DIR%\project-pkg\pkg-malloc
mklink /D %DIR%\project-app\app-alloc\vendor\alloc-one\vendor\pkg-malloc %DIR%\project-pkg\pkg-malloc
mklink /D %DIR%\project-app\app-alloc\vendor\alloc-two\vendor\pkg-malloc %DIR%\project-pkg\pkg-malloc
