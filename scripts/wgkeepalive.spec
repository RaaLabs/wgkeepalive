# The name should reflect the name part of the tar.gz archive
Name:           wgkeepalive
# The version should reflect the version part of the tar.gz arcive
Version:        0.1.0
Release:        0
Summary:        Package for wgkeepalive

License:        Apache 2.0
URL:            https://github.com/RaaLabs/wgkeepalive

# The name of the tar.gz archive as present in the SOURCE folder
Source0:        %{name}-%{version}-x86_64.tar.gz

Requires:       bash

# options x86_64 or noarch
BuildArch:      x86_64

%description
%{name}

# Unpack and organize the buildroot folder. Setup will basically unpack the tar.gz archive into the buildroot folder. The default buildroot folder will be named the same as the package, but since package name in this example is sftp and the folder structure when the tar.gz file here will be ./sftpproxy/... we need to specify what the name of the buildroot is, hence the use of -n below.
%prep
%setup -q

# This is typically where we do a make. Since we only have a binary to install we comment out this section.
%build

# Make the folder structure you want as it should be when installed on the client system, and also prepare/copy files needed within the structure. Bash scripting can also be used within this section.
%install
install -d -m 0755 %{buildroot}/usr/local/%{name}
install -m 0755 %{name} %{buildroot}/usr/local/%{name}/%{name}
install -m 0755 scripts/run.sh %{buildroot}/usr/local/%{name}/run.sh

mkdir -p %{buildroot}/usr/lib/systemd/system/
install -m 0755 scripts/%{name}.service %{buildroot}/usr/lib/systemd/system/%{name}.service
mkdir -p %{buildroot}/usr/lib/systemd/system/multi-user.target.wants
ln -s /usr/lib/systemd/system/%{name}.service %{buildroot}/usr/lib/systemd/system/multi-user.target.wants/%{name}.service

mkdir -p %{buildroot}/usr/share/clr-service-restart
ln -sf /usr/lib/systemd/system/%{name}.service %{buildroot}/usr/share/clr-service-restart/%{name}.service

# Specify all files as present from the buildroot in the %_install section above to be copied/installed when this package is installed by the user. Only the files or folders specified will actually be installed on the target system.

%files
# %license LICENSE
/usr/local/%{name}/%{name}
#/usr/local/trafficmonitor/run.sh
/usr/share/clr-service-restart/%{name}.service
/usr/lib/systemd/system/%{name}.service
/usr/lib/systemd/system/multi-user.target.wants/%{name}.service

%changelog