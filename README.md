## Golang NamePicker CLI App

**Overview:**
The NamePicker CLI is a simple command-line application written in Go that randomly selects a name from a provided list. This tool can be useful for various purposes, such as picking a winner from a pool of entries, selecting a team member for a task, or generating random names for games.

**Features:**

- **Random Selection:** Randomly selects a name from a given list.
- **Input Options:** Accepts names via command-line arguments or from a file.
- **Output Formatting:** Displays the selected name clearly and can format the output (e.g., uppercase, lowercase).
- **Error Handling:** Provides informative error messages for invalid inputs or file read issues.
- **Customizable Seed:** Allows users to set a random seed for reproducible results.
- **Help Command:** Displays help information and usage instructions.

**Installation:**
To install the NamePicker CLI app, follow these steps:

1. Ensure you have Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).
2. Clone the repository:

   ```bash
   git clone https://github.com/aletomasella/namepicker-cli.git
   cd namepicker-cli
   ```

3. Build the application:

   ```bash
   Make build
   ```

**Usage:**
The NamePicker CLI can be used in several ways:

1. **Select a Name from Command-Line Arguments:**

   ```bash
   namepicker Alice Bob Charlie
   ```

2. **Select a Name from a File:**

   Create a file named `names.txt` with one name per line:

   ```
   Alice
   Bob
   Charlie
   ```

   Then run:

   ```bash
   namepicker -file names.txt
   ```

3. **Get Help:**

   For help and usage instructions:

   ```bash
   namepicker -help
   ```
