Name:           xilinx-container-runtime
Version:        0.0.1
Release:        1%{?dist}
Summary:        Xilinx container runtime

License:        ASL 2.0
URL:            https://xilinx.com
Source0:        xilinx-container-runtime
Source1:        config.toml


%description
Xilinx-container-runtime is an extension of runc, with modification to add xilinx devices before running containers.

%prep
cp %{SOURCE0} %{SOURCE1} .


%install
rm -rf %{buildroot}
mkdir -p %{buildroot}%{_bindir}
install -m 755 -t %{buildroot}%{_bindir} xilinx-container-runtime
mkdir -p %{buildroot}%{_sysconfdir}/%{name}
install -m 755 -t %{buildroot}%{_sysconfdir}/xilinx-container-runtime config.toml

%files
%{_bindir}/%{name}
%{_sysconfdir}/%{name}/config.toml



%changelog
