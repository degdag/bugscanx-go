<h1 align="center">BugScanX-Go: SNI Bug Host Scanner</h1>

<p align="center">
   <i><b>BugScanX-Go</b> is an tool designed for scanning SNI bug hosts using various methods. It is a refined fork of BugScanner-Go, offering enhanced features and improved performance.</i>
</p>

<div align="center">
   <a href="https://github.com/Ayanrajpoot10/BugScanX-Go/stargazers">
      <img src="https://img.shields.io/github/stars/Ayanrajpoot10/BugScanX-Go?style=for-the-badge&color=green" alt="Stars Badge">
   </a>
   <a href="https://t.me/BugscanX">
      <img src="https://img.shields.io/badge/Telegram-Join%20Group-0088cc?style=for-the-badge&logo=telegram" alt="Telegram">
   </a>
   <a href="https://t.me/BugscanxChat">
      <img src="https://img.shields.io/badge/Telegram%20Chat-Join%20Chat-4c6ef5?style=for-the-badge&logo=telegram" alt="Telegram Chat">
   </a>
</div>

##  Changelog

###  Newly Added Features

- **Removed 302 Response Handling**: Skips 302 responses that redirect to recharge portals.
- **Comprehensive Server Saving**: Saves all server results under the "Others" category.
- **Expanded HTTP Methods**: Added support for methods like GET, PATCH, and PUT (default: HEAD).
- **New PING Scan**: A new method to scan using TCP ping.
- **Additional Changes**: Various minor enhancements and optimizations.

Here's an enhanced version of your installation and usage guide with a cleaner and more visually appealing structure:  


#  **Installation Guide**

###  **Step 1: Install BugScanX-Go**  
Use the following command to install the latest version of BugScanX-Go:  
```bash
go install -v github.com/Ayanrajpoot10/bugscanx-go@latest
```


###  **Step 2: Add Go Bin to PATH**  
Ensure Go binaries are accessible from anywhere by adding them to your PATH:  

#### For **Bash**:
```bash
echo 'export PATH="$PATH:$HOME/go/bin"' >> $HOME/.bashrc && source $HOME/.bashrc
```

#### For **Zsh**:
```bash
echo 'export PATH="$PATH:$HOME/go/bin"' >> $HOME/.zshrc && source $HOME/.zshrc
```


#  **How to Use**

Start by accessing the help menu to explore BugScanX-Go's options:  
```bash
bugscanx-go --help
```

###  **Preparation Before Scanning**  

1. **Install Subfinder**  
   To gather subdomains, install Subfinder or a similar tool by following the instructions at the [Subfinder Repository](https://github.com/projectdiscovery/subfinder#installation).  

2. **Save Subdomains to a File**  
   Use Subfinder to scan a domain and save the output:  
   ```bash
   subfinder -d example.com -o example.com.lst
   ```


###  **Scanning Examples**  

#### Direct Scan  
Scan directly using a file of subdomains:  
```bash
bugscanx-go scan direct -f example.txt -o cf.txt
```

#### CDN SSL Scan 
Perform an SSL scan through a CDN:  
```bash
bugscanx-go scan cdn-ssl --proxy-filename cf.txt --target ws.example.com
```
*The target server must respond with a 101 status code.*

#### Server Name Indication (SNI) Scan  
Run an SNI scan with custom threads and timeout:  
```bash
bugscanx-go scan sni -f example.com.txt --threads 16 --timeout 8 --deep 3
```

#### Scan Ping 
Perform a ping scan and save results:  
```bash
bugscanx-go scan ping -f example.txt --threads 15 -o save.txt
```

#### Scan DNS
check Dns wheater is it used for slow dns
```bash
bugscanx-go scan dns -f exaple.txt -o save.txt
