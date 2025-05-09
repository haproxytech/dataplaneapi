# /services/haproxy/configuration/frontends/{parent_name}/ssl_front_uses/{index}
ssl_front_uses_path() {
    echo "/services/haproxy/configuration/frontends/${1:?}/ssl_front_uses${2:+/$2}"
}
