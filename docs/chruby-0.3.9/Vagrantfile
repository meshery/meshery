# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant::Config.run do |config|
  # vagrant box add debian-wheezy-amd64 https://dl.dropboxusercontent.com/u/67225617/lxc-vagrant/lxc-wheezy64-puppet3-2013-07-27.box
  config.vm.define :debian do |debian|
    debian.vm.box = 'debian-wheezy-amd64'
  end

  # vagrant box add ubuntu-12.04-amd64 http://cloud-images.ubuntu.com/precise/current/precise-server-cloudimg-vagrant-amd64-disk1.box
  config.vm.define :ubuntu do |ubuntu|
    ubuntu.vm.box = 'ubuntu-12.04-amd64'
  end

  # vagrant box add rhel-6-amd64 http://puppetlabs.s3.amazonaws.com/pub/rhel60_64.box
  config.vm.define :redhat do |redhat|
    redhat.vm.box = 'rhel-6-amd64'
  end

  # vagrant box add centos-6-amd64 http://puppetlabs.s3.amazonaws.com/pub/centos4_64.box
  config.vm.define :centos do |centos|
    centos.vm.box = 'centos-6-amd64'
  end

  # vagrant box add freebsd-9.1-amd64 https://github.com/downloads/xironix/freebsd-vagrant/freebsd_amd64_ufs.box
  config.vm.define :freebsd do |freebsd|
    freebsd.vm.box = 'freebsd-9.1-amd64'
  end

  # vagrant box add openbsd-5.2-amd64 https://dl.dropbox.com/s/5ietqc3thdholuh/openbsd-52-64.box
  config.vm.define :openbsd do |openbsd|
    openbsd.vm.box = 'openbsd-5.2-amd64'
  end
end
