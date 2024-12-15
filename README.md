<<<<<<< HEAD
<h1 align="center">🔍🐛 BugScanX-Go: Advanced SNI Bug Host Scanner 🐞🔍</h1>

<p align="center">
   <i><b>BugScanX-Go</b> is an advanced tool designed for scanning SNI bug hosts using various methods. It is a refined fork of BugScanner-Go, offering enhanced features and improved performance.</i>
</p>

<div align="center">
   <a href="https://github.com/Ayanrajpoot10/BugScanX-Go/stargazers" style="margin: 0 10px;">
      <img src="https://img.shields.io/github/stars/Ayanrajpoot10/BugScanX-Go?style=for-the-badge&color=green" alt="Stars Badge" style="border-radius: 8px;">
   </a>
   <a href="https://t.me/BugscanX" style="margin: 0 10px;">
      <img src="https://img.shields.io/badge/Telegram-Join%20Group-0088cc?style=for-the-badge&logo=telegram" alt="Telegram" style="border-radius: 8px;">
   </a>
   <a href="https://t.me/BugscanxChat" style="margin: 0 10px;">
      <img src="https://img.shields.io/badge/Telegram%20Chat-Join%20Chat-4c6ef5?style=for-the-badge&logo=telegram" alt="Telegram Chat" style="border-radius: 8px;">
   </a>
</div>
=======
# BugScanX-Go

**A Better Version of BugScanner-Go**

Welcome to BugScanX-Go, the latest iteration of the BugScanner-Go project. This tool is designed to streamline bug hunting and security scanning with enhanced features, improved functionality, and robust performance. Help us reach **20 stars** on this repository, and we’ll release the latest code for everyone to enjoy!
>>>>>>> add171570839436abcbd6af51a1dd2f0e00b887d

---

<<<<<<< HEAD
---

## 📜 Changelog

### ✨ Newly Added Features

- **Removed 302 Response Handling**: Skips 302 responses that redirect to recharge portals.
- **Comprehensive Server Saving**: Saves all server results under the "Others" category.
- **Expanded HTTP Methods**: Added support for methods like GET, PATCH, and PUT (default: HEAD).
- **Dual Scheme Scanning**: Use the `--both-schemes` flag to scan both HTTP and HTTPS simultaneously.
- **New PING Scan**: A faster scanning option compared to other methods.
- **UI Improvements**: Enhanced color scheme and refined user interface elements for better usability.
- **Additional Changes**: Various minor enhancements and optimizations.

Here's an enhanced version of your installation and usage guide with a cleaner and more visually appealing structure:  

---

# 🛠️ **Installation Guide**

### 🚩 **Step 1: Install BugScanX-Go**  
Use the following command to install the latest version of BugScanX-Go:  
```bash
go install -v github.com/Ayanrajpoot/bugscanx-go@latest
```

---

### 🚩 **Step 2: Add Go Bin to PATH**  
Ensure Go binaries are accessible from anywhere by adding them to your PATH:  

#### For **Bash**:
```bash
echo 'export PATH="$PATH:$HOME/go/bin"' >> $HOME/.bashrc && source $HOME/.bashrc
```

#### For **Zsh**:
```bash
echo 'export PATH="$PATH:$HOME/go/bin"' >> $HOME/.zshrc && source $HOME/.zshrc
```

---

# 🚀 **How to Use**

Start by accessing the help menu to explore BugScanX-Go's options:  
```bash
bugscanx-go --help
```

---

### 🔍 **Preparation Before Scanning**  

1. **Install Subfinder**  
   To gather subdomains, install Subfinder or a similar tool by following the instructions at the [Subfinder Repository](https://github.com/projectdiscovery/subfinder#installation).  

2. **Save Subdomains to a File**  
   Use Subfinder to scan a domain and save the output:  
   ```bash
   subfinder -d example.com -o example.com.lst
   ```

---

### 🕵️ **Scanning Examples**  

#### **1️⃣ Direct Scan**  
Scan directly using a file of subdomains:  
```bash
bugscanx-go scan direct -f example.txt -o cf.txt
```

#### **2️⃣ CDN SSL Scan**  
Perform an SSL scan through a CDN:  
```bash
bugscanx-go scan cdn-ssl --proxy-filename cf.txt --target ws.example.com
```
*The target server must respond with a 101 status code.*

#### **3️⃣ Server Name Indication (SNI) Scan**  
Run an SNI scan with custom threads and timeout:  
```bash
bugscanx-go scan sni -f example.com.txt --threads 16 --timeout 8 --deep 3
```

#### **4️⃣ Scan Ping**  
Perform a ping scan and save results:  
```bash
bugscanx-go scan ping -f example.txt --threads 15 -o save.txt
```

#### **5️⃣ Scan DNS**
check Dns wheater is it used for slow dns
```bash
bugscanx-go scan dns -f exaple.txt -o save.txt
```
---

## Stargazers over time
[![Stargazers over time](https://starchart.cc/Ayanrajpoot10/bugscanx-go.svg?variant=adaptive)](https://starchart.cc/Ayanrajpoot10/bugscanx-go)
=======
## 🚀 Sneak Peek of BugScanX-Go (formerly BugHunter-Go):

Exciting updates are here! We’ve been working hard to make BugScanX-Go the ultimate bug-hunting tool. Here’s what’s new:

- **302 Response Removed**: Cleaner output by skipping 302 responses that redirect to recharge portals.
- **All Server Types Saved**: Now saves all server results under the "Others" category, not just Cloudflare, CloudFront, and Akamai. Check out the [example output file](https://t.me/bugscanxchat/28847).
- **More HTTP Methods**: Added support for methods like `GET`, `PATCH`, and `PUT` (default: `HEAD`). You can switch using the `--method` flag.
- **Both HTTP and HTTPS**: Use the `--both-schemes` flag to scan both HTTP and HTTPS at the same time.
- **New PING Scan**: A quick scanning option that’s faster than other methods.

---

## 📝 Planned Release
>>>>>>> add171570839436abcbd6af51a1dd2f0e00b887d

BugScanX-Go will be released along with a detailed tutorial for installation and usage once we reach **500 subscribers** or this post gets **150 reactions**. Stay tuned!

<<<<<<< HEAD
=======
---

## 💡 Got Feedback?

We’re always looking to improve BugScanX-Go. If you have ideas, suggestions, or feedback, let us know by creating an issue or joining the conversation in our [Telegram group](https://t.me/bugscanxchat).

---

## ⭐ How You Can Help
>>>>>>> add171570839436abcbd6af51a1dd2f0e00b887d

Support the development by starring this repository. Help us reach **20 stars** to unlock the latest version for everyone!

<<<<<<< HEAD

=======
---

## 📖 Documentation

Comprehensive documentation will be available upon release, covering:

1. Installation instructions
2. Command-line options and flags
3. Example use cases

---

Stay connected and help shape the future of BugScanX-Go!
>>>>>>> add171570839436abcbd6af51a1dd2f0e00b887d
