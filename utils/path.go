/*
@Time : 2021/3/9 10:36 PM
@Author : charley
@File : path
*/
package utils

import "os"

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
