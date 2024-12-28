# Contributing to Abyssal Kraken

Thank you for considering contributing to Abyssal Kraken! We welcome contributions of all types, including bug fixes, new features, documentation, and examples. Please take a moment to review these guidelines before starting.

---

## 💡 How to Contribute

1. **Fork the Repository**:
    - Create your own copy of the repository using the `Fork` button on GitHub.

2. **Clone Your Fork Locally**:
    - Clone the forked repository to your local machine:
      ```bash
      git clone https://github.com/abyssal-kraken/abyssalkraken.git
      ```

3. **Create a Branch for Your Changes**:
    - Use a descriptive branch name based on what you're working on:
      ```bash
      git checkout -b feature/new-feature
      ```

4. **Make Your Changes**:
    - Implement your modifications, add features, and fix bugs in the codebase.
    - Ensure all code is clear, concise, and well-commented.

5. **Test Your Changes**:
    - Run the unit tests to ensure your changes do not break existing functionality:
      ```bash
      go test ./...
      ```

6. **Commit Your Changes**:
    - Follow good commit practices, writing clear and descriptive commit messages:
      ```bash
      git add .
      git commit -m "Add a new feature to handle custom domain events"
      ```

7. **Push Your Changes**:
    - Push your branch to your forked repository:
      ```bash
      git push origin feature/new-feature
      ```

8. **Open a Pull Request**:
    - Go to the original repository on GitHub and open a pull request.
    - Provide a descriptive title and summary of the changes you’ve made.

---

## 🚦 Code of Conduct

To maintain a welcoming and productive environment, all contributors are expected to adhere to the project’s [Code of Conduct](./CODE_OF_CONDUCT.md).

---

## 🛠️ Development Guidelines

1. **Code Style**:
    - Follow Go’s formatting conventions using `gofmt`.
    - Be consistent with the surrounding code.

2. **Dependencies**:
    - Use Go Modules (`go.mod`) for managing dependencies.
    - Avoid adding unnecessary dependencies to keep the library lightweight.

3. **Testing**:
    - Write unit tests for any new functionality.
    - Ensure that all tests pass before submitting your pull request.

4. **Documentation**:
    - Update the `README.md` or other relevant files for any user-facing changes.

---

## 📢 Getting Help

If you’re unsure how to proceed, feel free to ask questions in the project’s [Discussions](https://github.com/abyssal-krakens/abyssalkraken/discussions) or open an issue.

---

## 🌟 Recognition

We truly appreciate your effort and dedication to improving Abyssal Kraken. All contributors will be recognized in this repository!

Thank you for contributing! 💙
