package main

import (
    "strings"
    "encoding/base64"
    "encoding/binary"
)

func path_to_id(path string) []uint64 {
    path_segments := strings.Split(path, "/")
    id := make([]uint64, len(path_segments))

    for i, _ := range path_segments {
        buf, err := base64.URLEncoding.DecodeString(path_segments[i])
        if err != nil {
            return []uint64{}
        }
        j, n := binary.Uvarint(buf)
        if n == 0 {
            return []uint64{}
        }
        id[i] = j;
    }
    return id
}

func id_to_path(id []uint64) string {
    path_segments := make([]string, len(id))
    for i, _ := range id {
        buf := make([]byte, 8)
        binary.PutUvarint(buf, id[i])
        path_segments[i] = base64.URLEncoding.EncodeToString(buf)
    }
    path := strings.Join(path_segments, "/")
    return path
}
