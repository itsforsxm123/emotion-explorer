
## Getting Started

**Prerequisites:**

1.  **Go:** Version 1.18 or later recommended.
2.  **C Compiler:** Fyne requires a C compiler for CGO. Follow the Fyne documentation instructions for your OS: [https://developer.fyne.io/started/](https://developer.fyne.io/started/)
    *   **Windows:** Requires `gcc` via MinGW (e.g., installed via MSYS2). Ensure GCC is correctly added to your system PATH.
    *   **macOS:** Requires Xcode Command Line Tools.
    *   **Linux:** Requires `gcc` and relevant development packages (e.g., `build-essential`, `libgl1-mesa-dev`, `xorg-dev` on Debian/Ubuntu).

**Running the Application:**

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/itsforsxm123/emotion-explorer.git
    cd emotion-explorer
    ```
2.  **Run from source:**
    ```bash
    go run ./cmd/emotion-explorer/
    ```
    *(The first run might take a moment to download dependencies.)*

## Current Development Stage & Next Steps

The core hierarchical navigation is functional, and initial GUI modernization steps have been taken by implementing a card-based layout with color integration and refactoring the UI code.

**Next Steps / Current Focus:**

The focus remains on **improving the User Experience and GUI**. Potential next steps include:

1.  **GUI Polish:**
    *   Add icons to the emotion cards.
    *   Implement theme switching (Light/Dark).
    *   Refine layouts, spacing, and visual presentation further.
2.  **Implement Detail View:** Create a dedicated view to display when a leaf node emotion is selected, showing details like name, description (if added to data), color, and parent.
3.  **Add Breadcrumbs:** Implement a visual breadcrumb trail (e.g., "Home > Angry > Aggressive") to improve navigation context.
4.  **Data Completion:** Add missing color codes to `emotions.json` (e.g., for "Hostile", "Provoked").

## Future Plans / Roadmap

### Phase 1: UI/UX Enhancements & Detail View (Current/Near Term)

*(Contains the items listed under "Next Steps" above)*

### Phase 2: Emotion Tracking & Journaling (Longer Term Vision)

Transform the application from a passive explorer into a personal emotion tracking tool:

1.  **Tracking Mechanism:** Allow users to select an emotion (likely a leaf node) and record it with a timestamp.
2.  **Persistence:** Implement data storage for tracked entries (e.g., using a local file or ideally a simple database like SQLite).
3.  **Journaling:** Enable users to add optional text notes to their tracked emotion entries for context.
4.  **Review Interface:**
    *   Create views to review past entries (e.g., a simple list view).
    *   Potentially implement a calendar view to browse entries by date.
5.  **Visualization:** Add graphing capabilities to show emotion trends over time (daily, weekly, monthly).

## Contributing

Contributions are welcome! Please feel free to open an issue to discuss potential changes or submit a pull request.

## License

This project is licensed under the MIT License - see the details below.

Copyright (c) 2025 itsforsxm123

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.