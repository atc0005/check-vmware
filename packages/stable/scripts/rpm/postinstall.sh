#!/bin/bash

# Copyright 2023 Adam Chalkley
#
# https://github.com/atc0005/check-vmware
#
# Licensed under the MIT License. See LICENSE file in the project root for
# full license information.

project_org="atc0005"
project_shortname="check-vmware"

project_fq_name="${project_org}/${project_shortname}"
project_url_base="https://github.com/${project_org}"
project_repo="${project_url_base}/${project_shortname}"
project_releases="${project_repo}/releases"
project_issues="${project_repo}/issues"
project_discussions="${project_repo}/discussions"

plugin_name_suffix=""
plugin_path="/usr/lib64/nagios/plugins"

#
# Set required SELinux context to allow plugin use when SELinux is enabled.
#

# Make sure we can locate the selinuxenabled binary.
if [ -x "$(command -v selinuxenabled)" ]; then
    selinuxenabled

    if [ $? -ne 0 ]; then
        echo -e "\n[--] SELinux is not enabled, skipping application of contexts."
    else
        # SELinux is enabled. Set context.
        echo -e "\nApplying SELinux contexts on plugins ..."

        for plugin_name in \
            check_vmware_tools \
            check_vmware_vcpus \
            check_vmware_vhw \
            check_vmware_hs2ds2vms \
            check_vmware_datastore_space \
            check_vmware_datastore_performance \
            check_vmware_snapshots_age \
            check_vmware_snapshots_count \
            check_vmware_snapshots_size \
            check_vmware_rps_memory \
            check_vmware_host_memory \
            check_vmware_host_cpu \
            check_vmware_vm_power_uptime \
            check_vmware_disk_consolidation \
            check_vmware_question \
            check_vmware_alarms \
            check_vmware_vm_backup_via_ca
        do

            echo -e "\tApplying SELinux contexts on ${plugin_path}/${plugin_name}${plugin_name_suffix}"

            chcon \
                --verbose \
                -t nagios_unconfined_plugin_exec_t \
                -u system_u \
                -r object_r \
                "${plugin_path}/${plugin_name}${plugin_name_suffix}"

            if [ $? -eq 0 ]; then
                echo -e "\t[OK] Successfully applied SELinux contexts on ${plugin_path}/${plugin_name}${plugin_name_suffix}"
            else
                echo -e "\t[!!] Failed to set SELinux contexts on ${plugin_path}/${plugin_name}${plugin_name_suffix}"
            fi

        done
    fi

else
    echo "[!!] Error: Failed to locate selinuxenabled command." >&2
fi



echo
echo "Thank you for installing packages provided by the ${project_fq_name} project!"
echo
echo "Project resources:"
echo
echo "- Obtain latest release: ${project_releases}"
echo "- View/Ask questions: ${project_discussions}"
echo "- View/Open issues: ${project_issues}"
echo


# Copying contexts from existing plugin:
#
# [root@atc0005-rhel8-test ~]# ls -laZ /usr/lib64/nagios/plugins/check_by_ssh
# -rwxr-xr-x. root root system_u:object_r:nagios_unconfined_plugin_exec_t:s0 /usr/lib64/nagios/plugins/check_by_ssh
# [root@atc0005-rhel8-test ~]# sudo chcon --verbose --reference /usr/lib64/nagios/plugins/check_by_ssh /usr/lib64/nagios/plugins/check_cert

# Explicitly applying contexts:
#
# [root@atc0005-rhel8-test ~]# sudo chcon --verbose -t nagios_unconfined_plugin_exec_t -u system_u -r object_r /usr/lib64/nagios/plugins/check_cert
# changing security context of '/usr/lib64/nagios/plugins/check_cert'
#
# [root@atc0005-rhel8-test ~]# ls -laZ /usr/lib64/nagios/plugins/check_cert
# -rwxr-xr-x. 1 root root system_u:object_r:nagios_unconfined_plugin_exec_t:s0 6035136 Dec 12 11:14 /usr/lib64/nagios/plugins/check_cert
