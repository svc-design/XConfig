#!/bin/sh
# Example commands demonstrating builtin modules
craftweave ansible all -i ../../inventory -m shell -a 'hostname'
craftweave ansible all -i ../../inventory -m command -a '/bin/echo hi'
craftweave ansible all -i ../../inventory -m copy -a '../../motd.tmpl:/tmp/motd-copy'
