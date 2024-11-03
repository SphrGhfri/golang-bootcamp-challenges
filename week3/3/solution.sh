#!/bin/bash

groupadd shared

useradd -m -G shared user1
useradd -m -G shared user2

mkdir /shared_files

touch /shared_files/shared_file


chown -R user1:shared /shared_files

chgrp shared /shared_files
chmod g+rw  /shared_files

userdel -r user1
userdel -r user2
groupdel shared
rm -rf /shared_files
