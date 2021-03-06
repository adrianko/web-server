Exec { path => ['/bin/', '/sbin/' , '/usr/bin/', '/usr/sbin/'] }

exec { 
    'update':
        command => 'apt-get update --fix-missing -y';

    'upgrade':
        command => 'apt-get upgrade -y',
        require => Exec['update'];

    'add-go-root':
        command => 'echo "export GOPATH=/home/vagrant/testing/pkg" >> .bash_profile',
        provider => shell,
        cwd => '/home/vagrant';
        
    'retrieve-code':
        command => 'cp -r web-server/*.go testing/',
        cwd => '/home/vagrant',
        require => File['/home/vagrant/testing'];
}

package {
    ['nginx', 'apache2', 'golang', 'git']:
        ensure => present,
        require => Exec['upgrade'];
}

service {
    'nginx':
        ensure => running,
        require => Package['nginx'];

    'apache2':
        ensure => running,
        require => Package['apache2'];
}

file {
    '/etc/nginx/sites-available/default':
        ensure => present,
        owner => root,
        group => root,
        mode => 644,
        source => '/vagrant/config/nginx/default',
        notify => Service['nginx'],
        require => Package['nginx'];

    '/etc/maester-http':
        ensure => present,
        owner => root,
        group => root,
        mode => 644,
        source => '/vagrant/config/maester-http';
        
    '/home/vagrant/testing':
        ensure => directory,
        owner => vagrant,
        group => vagrant;
    
    '/home/vagrant/testing/pkg':
        ensure => directory,
        owner => vagrant,
        group => vagrant;

    '/home/vagrant/build.sh':
        ensure => present,
        owner => vagrant,
        group => vagrant,
        mode => 755,
        source => '/vagrant/build.sh';

    '/home/vagrant/run-servers.sh':
        ensure => present,
        owner => vagrant,
        group => vagrant,
        mode => 755,
        source => '/vagrant/run-servers.sh';
}
