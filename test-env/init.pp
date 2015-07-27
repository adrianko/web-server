Exec { path => ['/bin/', '/sbin/' , '/usr/bin/', '/usr/sbin/'] }

exec { 
    'update':
        command => 'apt-get update --fix-missing -y';

    'upgrade':
        command => 'apt-get upgrade -y',
        require => Exec['update'];
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
}
