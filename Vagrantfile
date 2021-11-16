# -*- mode: ruby -*-
# vi: set ft=ruby :

unless Vagrant.has_plugin?("vagrant-vbguest")
  puts 'Installing vagrant-vbguest Plugin...'
  system('vagrant plugin install vagrant-vbguest')
end
unless Vagrant.has_plugin?("vagrant-reload")
  puts 'Installing vagrant-reload Plugin...'
  system('vagrant plugin install vagrant-reload')
end


Vagrant.configure("2") do |config|

  config.ssh.insert_key = false
  config.ssh.forward_agent = true
  config.ssh.forward_x11 = true
  config.vbguest.auto_update = true
  config.vm.box_check_update = false
  config.vm.boot_timeout = 6000


  config.vm.define "ledger" do |ledger|
    ledger.vm.box = "ubuntu/focal64"
    ledger.vm.hostname = "ledger"

    ledger.vm.network :private_network, ip: "10.10.50.2"
    if false then #Vagrant::Util::Platform.windows? then
      ledger.vm.synced_folder "shared", "/home/vagrant/shared",
        id: "shared", owner: "vagrant", group: "vagrant",
        mount_options: ["dmode=775","fmode=764"]
      ledger.vm.synced_folder "certshare", "/home/vagrant/certshare",
        id: "certshare", owner: "vagrant", group: "vagrant",
        mount_options: ["dmode=775","fmode=764"]
      ledger.vm.synced_folder "chaincode", "/home/vagrant/chaincode",
        id: "chaincode", owner: "vagrant", group: "vagrant",
        mount_options: ["dmode=775","fmode=764"]
    else
      ledger.vm.synced_folder "shared", "/home/vagrant/shared"
      ledger.vm.synced_folder "certshare", "/home/vagrant/certshare"
      ledger.vm.synced_folder "chaincode", "/home/vagrant/chaincode"
    end
    ledger.vm.provider "virtualbox" do |vb|
      vb.name = "ledger"
      opts = ["modifyvm", :id, "--natdnshostresolver1", "on"]
      vb.customize opts
      vb.memory = "1024"
      vb.cpus = 2
    end
    ledger.vm.provision :shell, path: "clean-provision.sh"
    ledger.vm.provision :shell, path: "general-provision.sh"
    #ledger.vm.provision :reload
    ledger.vm.provision :shell, path: "caserver-provision.sh", privileged: false
    ledger.vm.provision :shell, path: "nodes-provision.sh", privileged: false
    ledger.vm.provision :shell, path: "channel-provision.sh", privileged: false

    ledger.vm.provision :shell, path: "create_mod_certs.sh", privileged: false
    ledger.vm.provision :shell, path: "chaincode_install.sh", privileged: false
  end

  config.vm.define "appServer" do |appServer|
    appServer.vm.box = "ubuntu/focal64"
    appServer.vm.hostname = "appServer"
    appServer.vm.network :private_network, ip: "10.10.50.5"
    if false then #Vagrant::Util::Platform.windows? then
      appServer.vm.synced_folder "appServer", "/home/vagrant/appServer",
        id: "appServer", owner: "vagrant", group: "vagrant",
        mount_options: ["dmode=775","fmode=764"]
      appServer.vm.synced_folder "certshare", "/home/vagrant/certshare",
        id: "certshare", owner: "vagrant", group: "vagrant",
        mount_options: ["dmode=775","fmode=764"]
    else
      appServer.vm.synced_folder "appServer", "/home/vagrant/appServer"
      appServer.vm.synced_folder "certshare", "/home/vagrant/certshare"
    end
    appServer.vm.provider "virtualbox" do |vb|
      vb.name = "appServer"
      opts = ["modifyvm", :id, "--natdnshostresolver1", "on"]
      vb.customize opts
      vb.memory = "2048"
      vb.cpus = 2
    end
    appServer.vm.provision :shell, path: "general-provision.sh"
    appServer.vm.provision :shell, path: "appserver-provision.sh"
  end
end
