{% for fn in interface %}
alias {{fn}}='./mx send stdio {{node}}.{{fn}}'

_mx_{{fn}}_complete(){
   local args cur
   cur="${COMP_WORDS[COMP_CWORD]}"
   args='-r {% for arg in interface[fn].args %}-{{arg}} {% endfor %}'
   if [[ ${cur:0:1} == '-' ]] || [[ $COMP_CWORD == '1' ]]; then COMPREPLY=( $(compgen -W "${args}" -- ${cur}) ); fi
   return 0
}
complete -F _mx_{{fn}}_complete {{fn}}
{% endfor %}
