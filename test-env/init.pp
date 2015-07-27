Exec { path => ['/bin/', '/sbin/' , '/usr/bin/', '/usr/sbin/'] }

exec { 
    'update':
        command => 'apt-get update --fix-missing -y';

    'upgrade':
        command => 'apt-get upgrade -y',
        require => Exec['update'];
}

package {
    ['nginx', 'apache2', 'golang']:
        ensure => present,
        require => Exec['upgrade'];
}
