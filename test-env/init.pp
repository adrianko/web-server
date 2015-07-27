Exec { path => ['/bin/', '/sbin/' , '/usr/bin/', '/usr/sbin/'] }

exec { 
    'update':
        command => 'apt-get update --fix-missing -y';

    'upgrade':
        command => 'apt-get upgrade -y',
        require => Exec['update'];

    'create-go-root':
        command => 'mkdir .go',
        cwd => '/home/vagrant';

    'add-go-root':
        command => 'echo "export GOPATH=$HOME/.go" >> .bash_profile',
        cwd => '/home/vagrant';
    
    'create-testing-folder':
        command => 'mkdir testing',
        cwd => '/home/vagrant';
        
    'retrieve-code':
        command => 'cp -r web-server/*.go testing/',
        cwd => '/home/vagrant',
        require => Exec['create-testing-folder'];

    'get-build-dependencies':
        command => 'go get',
        cwd => '/home/vagrant/testing',
        require => [ Package['golang'], Exec['retrieve-code'] ];

    'build':
        command => 'go build server.go',
        cwd => '/home/vagrant/testing',
        require => [ Exec['get-build-dependencies'], Package['golang'] ];
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
