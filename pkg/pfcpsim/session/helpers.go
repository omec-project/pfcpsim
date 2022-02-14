/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2022 Open Networking Foundation
 */

package session

func contains(slice []uint8, element uint8) bool {
	if len(slice) == 0 {
		return false
	}

	for _, e := range slice {
		if e == element {
			return true
		}
	}

	return false
}
