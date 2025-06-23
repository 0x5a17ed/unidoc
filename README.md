# UniDoc 🦋

**Beautiful Unicode Markdown Renderer**

Transform your plain Markdown documents into stunning Unicode text with mathematical symbols, elegant typography, and visual formatting. Perfect for plain text environments, and anywhere you need rich formatting but HTML or Markdown is unavailable.


## ✨ Features

UniDoc supports a wide range of formatting enhancing your markdown experience with:

- 🎨 **Bold Strong Formatting:**  
  Bold text using mathematical sans-serif (𝗧𝗵𝗶𝘀 𝗶𝘀 𝗯𝗼𝗹𝗱)
- 📝 **Multiple Italic Styles:**  
  Choose from markers, script, or slanted sans-serif (𝘛𝘩𝘪𝘴 𝘪𝘴 𝘪𝘵𝘢𝘭𝘪𝘤)
- 📊 **Fancy List Numbering:**  
  Multi-level Unicode bullets and numbering (① ② ③, 🅐 🅑 🅒, ⅰ ⅱ ⅲ)
- 📦 **Beautiful Code Blocks:**  
  Unicode box-drawing characters with proper content
- 💬 **Nested Blockquotes:**  
  Visual hierarchy with stacked `┃` symbols
- ➖ **Smart Dashes:**  
  Convert `--` to en dash (–) and `---` to em dash (—)
- 🔗 **Rich Links & Images:**  
  Emojis and proper URL formatting
- 📋 **Proper List Spacing:**  
  Handles tight and loose lists correctly
- 💔 **Hard Line Break Support:**  
  Respects double-space line breaks
- ⚡ **Fast & Reliable:**  
  Built on the proven `goldmark` Markdown parser, ensuring compatibility with existing markdown files


## 📦 Installation

```bash
go install github.com/0x5a17ed/unidoc/cmd/unidoc@latest
```


## 📚 Usage

### Basic Usage

```bash
# From stdin
echo '# Hello **World**' | unidoc

# From file
unidoc README.md

# Save to file
cat document.md | unidoc > output.txt
```

Optionally combine with the `wl-copy` command to copy output directly to clipboard:

```bash
echo '# Hello **World**' | unidoc | wl-copy
```

```bash
wl-paste | unidoc | wl-copy
```

### 🎭 Italic Styles

#### Markers
```
*italic text* -> *italic text* (no special formatting)
```

#### Script
```
*italic text* -> 𝒾𝓉𝒶𝓁𝒾𝒸 𝓉ℯ𝓍𝓉
```

#### Sans-Serif Italic (Default)
```
*italic text* -> 𝘪𝘵𝘢𝘭𝘪𝘤 𝘵𝘦𝘹𝘵
```


### 📋 List Examples

#### Unordered Lists
```markdown
* First level
  * Second level
    * Third level
```

**Output:**
```
• First level
  ◦ Second level
    ▪ Third level
```

#### Ordered Lists
```markdown
1. First item
2. Second item
   1. Nested item
   2. Another nested
      1. Deep nesting
```

**Output:**
```
① First item
② Second item
  ⑴ Nested item
  ⑵ Another nested
    🅐 Deep nesting
```


### 💬 Blockquote Examples

```markdown
> This is a blockquote
>> This is nested
>>> Triple nested
```

**Output:**
```
┃ This is a blockquote
┃ ┃ This is nested
┃ ┃ ┃ Triple nested
```


## 📦 Code Block Example

    ```go
    func main() {
        fmt.Println("Hello, World!")
    }
    ```

**Output:**
```
┌────────────────────────────────────────────────────────────────┐
│ func main() {                                                  │
│     fmt.Println("Hello, World!")                               │
│ }                                                              │
└────────────────────────────────────────────────────────────────┘
```


### ➖ Smart Dashes

```markdown
Date range: 2020--2024
This is important --- very important --- to remember.
```

**Output:**
```
Date range: 2020–2024
This is important — very important — to remember.
```


## 🔧 Command Line Options

```
Usage:
  unidoc [OPTION]... [FILE]

Options:
  -h, --help                       Show help information
      --italic-style italicStyle   style for italic text (default slanted-sans-serif)
      --strong-style strongStyle   style for strong text (default bold-sans-serif)

Italic Styles:
  plain                       use regular text, no special formatting
  markers                     use simple markers around italic text: *text*
  script                      use Unicode script letters: 𝒯𝒽𝒾𝓈 𝒾𝓈 𝒾𝓉𝒶𝓁𝒾𝒸
  slanted-sans-serif          use Unicode sans-serif italic: 𝘛𝘩𝘪𝘴 𝘪𝘴 𝘪𝘵𝘢𝘭𝘪𝘤

Strong Styles:
  plain                       use regular text, no special formatting
  markers                     use simple markers around strong text: **text**
  bold-sans-serif             use mathematical bold sans-serif: 𝗧𝗵𝗶𝘀 𝗶𝘀 𝗯𝗼𝗹𝗱
```


## 📄 Complete Example

**Input Markdown:**

    # Project Report
    
    This is a **comprehensive** analysis of our *findings* from 2020--2024.
    
    ## Key Points
    
    * **Performance:** Improved significantly
    * **Issues:** Resolved most problems
      * Critical bugs fixed
      * Minor improvements made
        * splines were reticulated
    
    > **Note:** This data is preliminary --- further analysis needed.
    
    ### Code provided
    
    ```go
    package main
    
    import "fmt"
    
    func analyze() {
        fmt.Println("success")
    }
    ```
    
    ## Conclusion
    
    The project shows excellent progress!


**UniDoc Output:**
```
█ 𝗣𝗿𝗼𝗷𝗲𝗰𝘁 𝗥𝗲𝗽𝗼𝗿𝘁
══════════════════════════════════════════════════

This is a 𝗰𝗼𝗺𝗽𝗿𝗲𝗵𝗲𝗻𝘀𝗶𝘃𝗲 analysis of our 𝘧𝘪𝘯𝘥𝘪𝘯𝘨𝘴 from 2020–2024.

▓▓ 𝗞𝗲𝘆 𝗣𝗼𝗶𝗻𝘁𝘀
──────────────────────────────────────────────────

• 𝗣𝗲𝗿𝗳𝗼𝗿𝗺𝗮𝗻𝗰𝗲: Improved significantly
• 𝗜𝘀𝘀𝘂𝗲𝘀: Resolved most problems
  ◦ Critical bugs fixed
  ◦ Minor improvements made
    ▪ splines were reticulated

┃ 𝗡𝗼𝘁𝗲: This data is preliminary — further analysis needed.

┌────────────────────────────────────────────────────────────────┐
│ package main                                                   │
│                                                                │
│ import "fmt"                                                   │
│                                                                │
│ func analyze() {                                               │
│     fmt.Println("success")                                     │
│ }                                                              │
└────────────────────────────────────────────────────────────────┘

▓▓ 𝗖𝗼𝗻𝗰𝗹𝘂𝘀𝗶𝗼𝗻
──────────────────────────────────────────────────

The project shows excellent progress!

```


## 🌟 Inspiration

**UniDoc** transforms the humble markdown format into beautiful Unicode art, making terminal-based documentation and plain text environments visually stunning while maintaining full compatibility with existing markdown workflows.

Perfect for:
- Plain text emails and chat
- Documentation systems
- Terminal applications
- CLI tools output
- Anywhere rich formatting is desired but HTML is not supported


## 📜 License

This project is licensed under the MIT License -- see the [LICENSE](LICENSE) file for details.

---

<p align="center"><span style="font-weight: bold;">Made with ❤️ for sharing beautiful documents over plain text fields 🌟</span></p>
