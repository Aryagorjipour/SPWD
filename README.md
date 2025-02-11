# spwd - Smart Password Generator

`spwd` is a simple command-line application written in Go to generate strong passwords based on various levels of complexity. It supports customizable password lengths and modes ranging from very weak to unbreakable.

---

## Features

- **Password Modes**: Generate passwords of varying complexity (from very weak to unbreakable).
- **Customizable Length**: Set the desired length of the password.
- **Clipboard Support**: Automatically copies the generated password to your clipboard.
- **Cross-platform**: Works on both Windows and Linux.

---

## Modes

- **`-vw`**: Very Weak (only numbers)
- **`-w`**: Weak (numbers + lowercase letters)
- **`-m`**: Medium (numbers + lowercase + uppercase letters)
- **`-s`**: Strong (numbers + lowercase + uppercase + special characters)
- **`-vs`**: Very Strong (numbers + lowercase + uppercase + special characters with more variety)
- **`-xb`**: Unbreakable (highly random characters from all types)

---

## Installation

### Windows

1. Download the `install_windows.bat` script.
2. Run the script by double-clicking it. This will install `spwd` to your system PATH, allowing you to use it from any terminal.
3. Make sure you have [Go](https://golang.org/dl/) installed before running the script.

### Linux

1. Download the `install_linux.sh` script.
2. Run the following commands to install the script:
   ```bash
   chmod +x install_linux.sh
   ./install_linux.sh
   ```
3. This will install `spwd` to `/usr/local/bin`, making it available globally in your terminal.

---

## Usage

Once the application is installed, you can run it directly from your terminal.

### Generating Passwords

To generate a password, use the following syntax:

```bash
spwd -L <length> -M <mode>
```

- **`-L <length>`**: The length of the password.
- **`-M <mode>`**: The complexity of the password (choose from `-vw`, `-w`, `-m`, `-s`, `-vs`, `-xb`).

### Example Commands

- **Very Weak (length 8)**:
  ```bash
  spwd -L 8 -M vw
  ```
  This will generate a very weak password of 8 characters, consisting of only digits.

- **Strong (length 16)**:
  ```bash
  spwd -L 16 -M s
  ```
  This will generate a strong password of 16 characters, consisting of numbers, lowercase letters, uppercase letters, and special characters.

- **Unbreakable (length 20)**:
  ```bash
  spwd -L 20 -M xb
  ```
  This will generate an unbreakable password of 20 characters.

---

## Clipboard Support

By default, the generated password will be copied to your clipboard after it is created. You can paste it directly into any form or password manager.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Contributing

If you’d like to contribute, please fork the repository and submit a pull request. 
Make sure to follow the Go code style.

---

## Acknowledgments

- Thanks to the Go community for making such a powerful and easy-to-use programming language.
- Clipboard functionality provided by the [atotto/clipboard](https://github.com/atotto/clipboard) package.
