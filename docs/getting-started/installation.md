# Getting Started with Dokku

## What is Dokku?

Dokku is an extensible, open source Platform as a Service that runs on a single server of your choice.

To start using Dokku, you'll need a system that meets the following minimum requirements:

- A fresh installation of [Ubuntu 16.04 x64](https://www.ubuntu.com/download), [Ubuntu 14.04 x64](https://www.ubuntu.com/download), [Debian 8.2 x64](https://www.debian.org/distrib/) or [CentOS 7 x64](https://www.centos.org/download/) *(experimental)* with the FQDN set <sup>[1]</sup>
- At least `1GB` of system memory <sup>[2]</sup>

You can *optionally* have a domain name pointed at the host's IP, though this is not necessary.

Dokku is designed for usage on a fresh VM installation, and should install all necessary dependencies if installing via the bootstrap method.

### Installing the latest stable version

#### 1. Install dokku

To install the latest stable version of dokku, you can run the following shell commands:

```shell
# for debian systems, installs Dokku via apt-get
wget https://raw.githubusercontent.com/dokku/dokku/v0.8.0/bootstrap.sh;
sudo DOKKU_TAG=v0.8.0 bash bootstrap.sh
```

The installation process takes about 5-10 minutes, depending upon internet connection speed.

#### 2. Setup SSH key and Virtualhost Settings

Once the installation is complete, you can open a browser to setup your SSH key and virtualhost settings. Open your browser of choice and navigate to the host's IP address - or the domain you assigned to that IP previously - and configure Dokku via the web admin.

>**Warning:** If you don't complete setup via the web installer (even if you set up SSH keys and virtual hosts otherwise) your Dokku installation will remain vulnerable to anyone finding the setup page and inserting their key.

>**Warning:** Web installer is not available on CentOS and Arch Linux. You will need to configure [SSH keys](/dokku/deployment/user-management/#adding-ssh-keys) and [virtual hosts](/dokku/configuration/domains/#customizing-hostnames) using dokku command line interface.

#### 3. Deploy your first application
Once you save your settings, the web admin will self-terminate and you should be able to run or deploy to the Dokku installation.

### Installing via other methods

For various reasons, certain hosting providers may have other steps that should be preferred to the above. If hosted on any of the following popular hosts, please follow the linked to instructions:

- [Digital Ocean Installation Notes](/dokku/getting-started/install/digitalocean)
- [DreamHost Cloud Installation Notes](/dokku/getting-started/install/dreamhost/)
- [Linode Installation Notes](/dokku/getting-started/install/linode/)
- [Microsoft Azure Installation Notes](/dokku/getting-started/install/azure/)

As well, you may wish to customize your installation in some other fashion. or experiment with vagrant. The guides below should get you started:

- [Debian Package Installation Notes](/dokku/getting-started/install/debian/)
- [RPM Package Installation Notes](/dokku/getting-started/install/rpm/)
- [Vagrant Installation Notes](/dokku/getting-started/install/vagrant/)
- [Advanced Install Customization](/dokku/getting-started/advanced-installation/)

---

- <sup>[1]: To check whether your system has an fqdn set, run `sudo hostname -f`</sup>
- <sup>[2]: If your system has less than 1GB of memory, you can use [this workaround](/dokku/getting-started/advanced-installation/#vms-with-less-than-1gb-of-memory).</sup>
