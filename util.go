package batis

import "regexp"

var tablePattern = regexp.MustCompile(`(?i)\s+from\s+[\s\S]*[where|group by]{0,1}`)

func tableList(sql string) []string  {
	return tablePattern.FindAllString(sql, -1)
}



