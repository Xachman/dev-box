#!/bin/bash
projects_path="/home/xach/dev-box-projects"
prefix="xachs_"
workspace_name="servicemaster"
image=$1
path=$2

function start_container {
    echo "docker run --name $prefix$workspace_name -v $projects_path$path:/projects$path $image"
}

start_container