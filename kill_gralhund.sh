ps | grep "./[d]ist/main -t" | awk '
{ 
  if($1!="") {
      print "killing gralhund: "$1
      system("kill " $1)
  }
}' 