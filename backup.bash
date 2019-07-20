#备份重要文件
filebak="/etc/system_bak"

if [ -e "$filebak" ];then
filebak_panduan=0
else
mkdir /etc/system_bak
chown root /etc/system_bak
chmod 750 /etc/system_bak
filebak_panduan=1
fi

cp /etc/login.defs /etc/system_bak/login.defs.bak
cp /etc/security/limits.conf /etc/system_bak/security_limits.conf.bak
cp /etc/profile /etc/system_bak/profile.bak
cp /etc/pam.d/system-auth /etc/system_bak/pam.d_system-auth.bak
cp /etc/inittab /etc/system_bak/inittab.bak
cp /etc/motd /etc/system_bak/motd.bak
cp /etc/xinetd.conf /etc/system_bak/xinetd.conf.bak
cp /etc/group /etc/system_bak/group.bak
cp /etc/shadow /etc/system_bak/shadow.bak
cp /etc/services /etc/system_bak/services.bak
cp /etc/passwd /etc/system_bak/passwd.bak
if [ -f "/etc/grub.conf" ];then
 cp /etc/grub.conf /etc/system_bak/grub.conf.bak
 cp /boot/grub/grub.conf /etc/system_bak/boot_grub_grub.conf.bak
fi
if [ -f "/etc/lilo.conf" ]; then
  cp /etc/lilo.conf /etc/system_bak/lilo.conf.bak
fi
cp /etc/ssh_banner /etc/system_bak/ssh_banner.bak
cp /etc/ssh/sshd_config /etc/system_bak/sshd_config.bak
cp /etc/aliases /etc/system_bak/aliases.bak