#!/bin/sh
find `pwd` -name "*.sln" -exec dotnet restore "{}" \; -exec nuget restore "{}" \;
