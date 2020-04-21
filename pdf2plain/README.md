## pdf2plain

Can be used to convert pdf to plain text for downstream process of bioextr. 


## Installation

Xpdf command line tools (required):

- Linux 32/64-bit: [download](https://xpdfreader-dl.s3.amazonaws.com/xpdf-tools-linux-4.02.tar.gz) ([GPG signature](https://xpdfreader-dl.s3.amazonaws.com/xpdf-tools-linux-4.02.tar.gz.sig))
- Windows 32/64-bit: [download](https://xpdfreader-dl.s3.amazonaws.com/xpdf-tools-win-4.02.zip) ([GPG signature](https://xpdfreader-dl.s3.amazonaws.com/xpdf-tools-win-4.02.zip.sig))
- Mac 64-bit: [download](https://xpdfreader-dl.s3.amazonaws.com/xpdf-tools-mac-4.02.tar.gz) ([GPG signature](https://xpdfreader-dl.s3.amazonaws.com/xpdf-tools-mac-4.02.tar.gz.sig))

```bash
# windows
wget https://github.com/openanno/bioextr/releases/download/v0.1.0/pdf2plain.exe

# osx
wget https://github.com/openanno/bioextr/releases/download/v0.1.0/pdf2plain_osx
mv pdf2plain_osx pdf2plain
chmod a+x pdf2plain

# linux
wget https://github.com/openanno/bioextr/releases/download/v0.1.0/pdf2plain_linux64
mv pdf2plain_linux64 pdf2plain
chmod a+x pdf2plain

# get latest version
go get -u github.com/openanno/bioextr/pdf2plain
```

## Usage

```
pdf2plain _examples/Multi-omic_approaches_to_improve_outcome_for_T-cel.pdf -o out.text
```

## Maintainer

- [@Jianfeng](https://github.com/Miachol)

## License

Academic Free License version 3.0

