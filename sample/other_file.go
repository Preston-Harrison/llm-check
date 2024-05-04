package sample

func writeToFile(data []byte) error {
    file, err := openFileForWriting()
    if err != nil {
        return err
    }

    _, err = file.Write([]byte("hello world"))
    if err != nil {
        return err
    }

    return nil
}
