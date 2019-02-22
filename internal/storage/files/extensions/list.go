package extensions

import (
	. "github.com/tags-drive/core/cmd"
)

// List with all supported types of files
var extensionsList = []Ext{
	// Archives
	{Ext: ".7z", FileType: FileTypeArchive, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".pkg", FileType: FileTypeArchive, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".rar", FileType: FileTypeArchive, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".zip", FileType: FileTypeArchive, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".tar.gz", FileType: FileTypeArchive, Supported: false, PreviewType: TypeUnsupported},

	// Audio
	{Ext: ".aif", FileType: FileTypeAudio, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".mpa", FileType: FileTypeAudio, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".wma", FileType: FileTypeAudio, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".wpl", FileType: FileTypeAudio, Supported: false, PreviewType: TypeUnsupported},
	// supproted
	{Ext: ".ogg", FileType: FileTypeAudio, Supported: true, PreviewType: MediaTypeAudioOGG},
	{Ext: ".wav", FileType: FileTypeAudio, Supported: true, PreviewType: MediaTypeAudioWAV},
	{Ext: ".mp3", FileType: FileTypeAudio, Supported: true, PreviewType: MediaTypeAudioMP3},

	// Image
	{Ext: ".bmp", FileType: FileTypeImage, Supported: true, PreviewType: FileTypeImage},
	{Ext: ".gif", FileType: FileTypeImage, Supported: true, PreviewType: FileTypeImage},
	{Ext: ".ico", FileType: FileTypeImage, Supported: true, PreviewType: FileTypeImage},
	{Ext: ".jpg", FileType: FileTypeImage, Supported: true, PreviewType: FileTypeImage},
	{Ext: ".png", FileType: FileTypeImage, Supported: true, PreviewType: FileTypeImage},
	{Ext: ".svg", FileType: FileTypeImage, Supported: true, PreviewType: FileTypeImage},
	{Ext: ".jpeg", FileType: FileTypeImage, Supported: true, PreviewType: FileTypeImage},

	// Video
	{Ext: ".avi", FileType: FileTypeVideo, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".mkv", FileType: FileTypeVideo, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".mov", FileType: FileTypeVideo, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".mpg", FileType: FileTypeVideo, Supported: false, PreviewType: TypeUnsupported},
	{Ext: ".mpeg", FileType: FileTypeVideo, Supported: false, PreviewType: TypeUnsupported},
	// supported
	{Ext: ".mp4", FileType: FileTypeVideo, Supported: true, PreviewType: MediaTypeVideoMP4},
	{Ext: ".webm", FileType: FileTypeVideo, Supported: true, PreviewType: MediaTypeVideoWebM},

	// Text (from https://github.com/github/linguist/blob/master/lib/linguist/languages.yml)
	{Ext: ".cfg", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // HAProxy
	{Ext: ".m", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Mercury
	{Ext: ".pic", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Pic
	{Ext: ".pb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // PureBasic
	{Ext: ".mm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Objective-C++
	{Ext: ".asn", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // ASN.1
	{Ext: ".d", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // DTrace
	{Ext: ".self", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Self
	{Ext: ".marko", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Marko
	{Ext: ".lean", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Lean
	{Ext: ".nl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // NewLisp
	{Ext: ".conllu", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // CoNLL-U
	{Ext: ".reb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Rebol
	{Ext: ".robot", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // RobotFramework
	{Ext: ".xpm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // X PixMap
	{Ext: ".properties", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},    // Java Properties
	{Ext: ".pas", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Pascal
	{Ext: ".brs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Brightscript
	{Ext: ".owl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Web Ontology Language
	{Ext: ".q", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // HiveQL
	{Ext: ".cshtml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // HTML+Razor
	{Ext: ".e", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Eiffel
	{Ext: ".befunge", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},       // Befunge
	{Ext: ".raml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // RAML
	{Ext: ".opa", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Opa
	{Ext: ".ex", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Elixir
	{Ext: ".ur", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // UrWeb
	{Ext: ".mq4", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // MQL4
	{Ext: ".mq5", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // MQL5
	{Ext: ".agda", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Agda
	{Ext: ".thrift", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Thrift
	{Ext: ".xm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Logos
	{Ext: ".srt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // SRecode Template
	{Ext: ".jinja", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // HTML+Django
	{Ext: ".sass", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Sass
	{Ext: ".pbt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // PowerBuilder
	{Ext: ".edn", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // edn
	{Ext: ".sublime-build", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText}, // JSON with Comments
	{Ext: ".blade", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Blade
	{Ext: ".pp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Puppet
	{Ext: ".cr", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Crystal
	{Ext: ".m", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Objective-C
	{Ext: ".bison", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Bison
	{Ext: ".druby", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Mirah
	{Ext: ".gn", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // GN
	{Ext: ".ebuild", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Gentoo Ebuild
	{Ext: ".j", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Objective-J
	{Ext: ".bdf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Glyph Bitmap Distribution Format
	{Ext: ".ftl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // FreeMarker
	{Ext: ".f90", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Fortran
	{Ext: ".swift", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Swift
	{Ext: ".phtml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // HTML+PHP
	{Ext: ".ipynb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Jupyter Notebook
	{Ext: ".tpl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Smarty
	{Ext: ".cp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Component Pascal
	{Ext: ".cob", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // COBOL
	{Ext: ".go", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Go
	{Ext: ".ring", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Ring
	{Ext: ".php", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // PHP
	{Ext: ".csd", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Csound Document
	{Ext: ".bsl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // 1C Enterprise
	{Ext: ".i3", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Modula-3
	{Ext: ".http", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // HTTP
	{Ext: ".ms", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // MAXScript
	{Ext: ".raw", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Raw token data
	{Ext: ".mak", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Makefile
	{Ext: ".bb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // BitBake
	{Ext: ".ml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // OCaml
	{Ext: ".textile", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},       // Textile
	{Ext: ".m4", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // M4
	{Ext: ".sed", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // sed
	{Ext: ".v", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Verilog
	{Ext: ".do", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Stata
	{Ext: ".gitignore", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // Ignore List
	{Ext: ".red", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Red
	{Ext: ".ahk", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // AutoHotkey
	{Ext: ".yar", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // YARA
	{Ext: ".d-objdump", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // D-ObjDump
	{Ext: ".xpl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // XProc
	{Ext: ".xsp-config", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},    // XPages
	{Ext: ".volt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Volt
	{Ext: ".as", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // AngelScript
	{Ext: ".lsl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // LSL
	{Ext: ".jade", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Pug
	{Ext: ".rs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // RenderScript
	{Ext: ".minid", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // MiniD
	{Ext: ".t", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Terra
	{Ext: ".coffee", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // CoffeeScript
	{Ext: ".html", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // HTML
	{Ext: ".awk", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Awk
	{Ext: ".dot", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Graphviz (DOT)
	{Ext: ".org", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Org
	{Ext: ".idr", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Idris
	{Ext: ".pcss", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // PostCSS
	{Ext: ".pgsql", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // PLpgSQL
	{Ext: ".agc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Apollo Guidance Computer
	{Ext: ".mod", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Modula-2
	{Ext: ".c", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // C
	{Ext: ".myt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Myghty
	{Ext: ".hlsl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // HLSL
	{Ext: ".p4", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // P4
	{Ext: ".thy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Isabelle
	{Ext: ".eex", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // HTML+EEX
	{Ext: ".scaml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Scaml
	{Ext: ".clw", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Clarion
	{Ext: ".sh-session", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},    // ShellSession
	{Ext: ".pod", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Pod 6
	{Ext: ".graphql", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},       // GraphQL
	{Ext: ".lhs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Literate Haskell
	{Ext: ".erb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // HTML+ERB
	{Ext: ".au3", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // AutoIt
	{Ext: ".gs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Genie
	{Ext: ".moo", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Moocode
	{Ext: ".hcl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // HCL
	{Ext: ".chs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // C2hs Haskell
	{Ext: ".jsp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Java Server Pages
	{Ext: ".sp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // SourcePawn
	{Ext: ".cmake", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // CMake
	{Ext: ".clj", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Clojure
	{Ext: ".shen", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Shen
	{Ext: ".wdl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // wdl
	{Ext: ".hxml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // HXML
	{Ext: ".yml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // YAML
	{Ext: ".tcl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Tcl
	{Ext: ".nf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Nextflow
	{Ext: ".ini", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // INI
	{Ext: ".xtend", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Xtend
	{Ext: ".pytb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Python traceback
	{Ext: ".epj", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Ecere Projects
	{Ext: ".glsl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // GLSL
	{Ext: ".uc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // UnrealScript
	{Ext: ".eb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Easybuild
	{Ext: ".jq", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // JSONiq
	{Ext: ".vhdl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // VHDL
	{Ext: ".pyx", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Cython
	{Ext: ".sv", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // SystemVerilog
	{Ext: ".objdump", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},       // ObjDump
	{Ext: ".smt2", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // SMT
	{Ext: ".regexp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Regular Expression
	{Ext: ".gitconfig", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // Git Config
	{Ext: ".json5", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // JSON5
	{Ext: ".yasnippet", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // YASnippet
	{Ext: ".anim", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Unity3D Asset
	{Ext: ".b", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Limbo
	{Ext: ".rexx", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // REXX
	{Ext: ".lfe", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // LFE
	{Ext: ".oxygene", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},       // Oxygene
	{Ext: ".bb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // BlitzBasic
	{Ext: ".asm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Assembly
	{Ext: ".flex", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // JFlex
	{Ext: ".cs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // C#
	{Ext: ".t", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Turing
	{Ext: ".pde", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Processing
	{Ext: ".fx", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // FLUX
	{Ext: ".webidl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // WebIDL
	{Ext: ".json", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // JSON
	{Ext: ".asc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Public Key
	{Ext: ".matlab", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // MATLAB
	{Ext: ".ls", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // LiveScript
	{Ext: ".stan", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Stan
	{Ext: ".mtml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // MTML
	{Ext: ".pl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Prolog
	{Ext: ".jsx", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // JSX
	{Ext: ".proto", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Protocol Buffer
	{Ext: ".ch", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Charity
	{Ext: ".xml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // XML
	{Ext: ".svg", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // SVG
	{Ext: ".mathematica", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},   // Mathematica
	{Ext: ".lol", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // LOLCODE
	{Ext: ".ooc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // ooc
	{Ext: ".xslt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // XSLT
	{Ext: ".rmd", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // RMarkdown
	{Ext: ".afm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Adobe Font Metrics
	{Ext: ".mask", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Mask
	{Ext: ".sch", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // KiCad Schematic
	{Ext: ".em", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // EmberScript
	{Ext: ".lvproj", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // LabVIEW
	{Ext: ".n", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Nemerle
	{Ext: ".cu", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Cuda
	{Ext: ".krl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // KRL
	{Ext: ".vim", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Vim script
	{Ext: ".pony", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Pony
	{Ext: ".sci", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Scilab
	{Ext: ".1", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Roff Manpage
	{Ext: ".rpy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Ren'Py
	{Ext: ".sparql", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // SPARQL
	{Ext: ".applescript", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},   // AppleScript
	{Ext: ".txt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Text
	{Ext: ".sage", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Sage
	{Ext: ".ck", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // ChucK
	{Ext: ".g4", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // ANTLR
	{Ext: ".fs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // F#
	{Ext: ".ls", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // LoomScript
	{Ext: ".fy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Fancy
	{Ext: ".fst", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // F*
	{Ext: ".pir", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Parrot Internal Representation
	{Ext: ".st", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Smalltalk
	{Ext: ".ice", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Slice
	{Ext: ".monkey", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Monkey
	{Ext: ".pogo", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // PogoScript
	{Ext: ".el", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Emacs Lisp
	{Ext: ".js", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // JavaScript
	{Ext: ".pro", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // QMake
	{Ext: ".rs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Rust
	{Ext: ".abap", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // ABAP
	{Ext: ".pasm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Parrot Assembly
	{Ext: ".cw", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Redcode
	{Ext: ".sl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Slash
	{Ext: ".l", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // PicoLisp
	{Ext: ".spec", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // RPM Spec
	{Ext: ".erl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Erlang
	{Ext: ".mms", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Module Management System
	{Ext: ".dae", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // COLLADA
	{Ext: ".scm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Scheme
	{Ext: ".nut", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Squirrel
	{Ext: ".py", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Python
	{Ext: ".nanorc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // nanorc
	{Ext: ".latte", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Latte
	{Ext: ".ne", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Nearley
	{Ext: ".iss", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Inno Setup
	{Ext: ".ebnf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // EBNF
	{Ext: ".ipf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // IGOR Pro
	{Ext: ".chpl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Chapel
	{Ext: ".coq", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Coq
	{Ext: ".dylan", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Dylan
	{Ext: ".lagda", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Literate Agda
	{Ext: ".sch", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Eagle
	{Ext: ".gradle", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Gradle
	{Ext: ".clp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // CLIPS
	{Ext: ".axs.erb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},       // NetLinx+ERB
	{Ext: ".eclass", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Gentoo Eclass
	{Ext: ".xbm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // X BitMap
	{Ext: ".als", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Alloy
	{Ext: ".groovy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Groovy
	{Ext: ".w", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // CWeb
	{Ext: ".ol", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Jolie
	{Ext: ".pls", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // PLSQL
	{Ext: ".purs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // PureScript
	{Ext: ".jl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Julia
	{Ext: ".bf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // HyPhy
	{Ext: ".q", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // q
	{Ext: ".hs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Haskell
	{Ext: ".ncl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // NCL
	{Ext: ".vb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Visual Basic
	{Ext: ".io", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Io
	{Ext: ".rg", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Rouge
	{Ext: ".haml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Haml
	{Ext: ".djs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Dogescript
	{Ext: ".ps1", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // PowerShell
	{Ext: ".ts", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // TypeScript
	{Ext: ".dart", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Dart
	{Ext: ".edc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Edje Data Collection
	{Ext: ".vcl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // VCL
	{Ext: ".zig", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Zig
	{Ext: ".ceylon", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Ceylon
	{Ext: ".fr", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Frege
	{Ext: ".pro", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // IDL
	{Ext: ".g", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // GAP
	{Ext: ".aj", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // AspectJ
	{Ext: ".sh", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Shell
	{Ext: ".orc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Csound
	{Ext: ".tcsh", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Tcsh
	{Ext: ".prg", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // xBase
	{Ext: ".elm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Elm
	{Ext: ".jison", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Jison
	{Ext: ".x", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // RPC
	{Ext: ".desktop", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},       // desktop
	{Ext: ".sc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // SuperCollider
	{Ext: ".nginxconf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // Nginx
	{Ext: ".re", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Reason
	{Ext: ".yang", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // YANG
	{Ext: ".com", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // DIGITAL Command Language
	{Ext: ".sas", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // SAS
	{Ext: ".ninja", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Ninja
	{Ext: ".grace", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Grace
	{Ext: ".cl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // OpenCL
	{Ext: ".d", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // D
	{Ext: ".creole", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Creole
	{Ext: ".kt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Kotlin
	{Ext: ".opal", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Opal
	{Ext: ".8xp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // TI Program
	{Ext: ".ML", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Standard ML
	{Ext: ".cfc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // ColdFusion CFC
	{Ext: ".bat", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Batchfile
	{Ext: ".oz", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Oz
	{Ext: ".ox", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Ox
	{Ext: ".gsp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Groovy Server Pages
	{Ext: ".roff", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Roff
	{Ext: ".rl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Ragel
	{Ext: ".gs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Gosu
	{Ext: ".handlebars", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},    // Handlebars
	{Ext: ".less", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Less
	{Ext: ".zone", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // DNS Zone
	{Ext: ".pd", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Pure Data
	{Ext: ".ecr", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // HTML+ECR
	{Ext: ".kicad_pcb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // KiCad Layout
	{Ext: ".ld", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Linker Script
	{Ext: ".b", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Brainfuck
	{Ext: ".f", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Filebench WML
	{Ext: ".apl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // APL
	{Ext: ".hh", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Hack
	{Ext: ".toc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // World of Warcraft Addon Data
	{Ext: ".numpy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // NumPy
	{Ext: ".sqf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // SQF
	{Ext: ".glf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Glyph
	{Ext: ".fea", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // OpenType Feature File
	{Ext: ".cy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Cycript
	{Ext: ".java", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Java
	{Ext: ".scala", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Scala
	{Ext: ".scad", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // OpenSCAD
	{Ext: ".apacheconf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},    // ApacheConf
	{Ext: ".pl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Perl
	{Ext: ".asy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // LTspice Symbol
	{Ext: ".mediawiki", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // MediaWiki
	{Ext: ".vue", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Vue
	{Ext: ".gd", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // GDScript
	{Ext: ".gbr", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Gerber Image
	{Ext: ".capnp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Cap'n Proto
	{Ext: ".factor", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Factor
	{Ext: ".reg", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Windows Registry Entries
	{Ext: ".darcspatch", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},    // Darcs Patch
	{Ext: ".fth", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Forth
	{Ext: ".asy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Asymptote
	{Ext: ".hy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Hy
	{Ext: ".j", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Jasmin
	{Ext: ".ec", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // eC
	{Ext: ".scss", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // SCSS
	{Ext: ".cls", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Apex
	{Ext: ".l", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Lex
	{Ext: ".rb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Ruby
	{Ext: ".ly", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // LilyPond
	{Ext: ".cl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Cool
	{Ext: ".zimpl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Zimpl
	{Ext: ".kid", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Genshi
	{Ext: ".golo", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Golo
	{Ext: ".cson", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // CSON
	{Ext: ".sql", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // SQL
	{Ext: ".metal", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Metal
	{Ext: ".gml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Graph Modeling Language
	{Ext: ".md", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // GCC Machine Description
	{Ext: ".ni", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Inform 7
	{Ext: ".lgt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Logtalk
	{Ext: ".mo", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Modelica
	{Ext: ".m4", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // M4Sugar
	{Ext: ".boo", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Boo
	{Ext: ".csv", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // CSV
	{Ext: ".eq", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // EQ
	{Ext: ".mtl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Wavefront Material
	{Ext: ".css", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // CSS
	{Ext: ".uno", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Uno
	{Ext: ".ttl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Turtle
	{Ext: ".c-objdump", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // C-ObjDump
	{Ext: ".rdoc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // RDoc
	{Ext: ".abnf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // ABNF
	{Ext: ".ampl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // AMPL
	{Ext: ".cfm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // ColdFusion
	{Ext: ".cirru", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Cirru
	{Ext: ".rst", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // reStructuredText
	{Ext: ".hb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Harbour
	{Ext: ".y", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Yacc
	{Ext: ".g", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // G-code
	{Ext: ".xojo_code", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // Xojo
	{Ext: ".srt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // SubRip Text
	{Ext: ".bmx", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // BlitzMax
	{Ext: ".pig", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // PigLatin
	{Ext: ".tl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Type Language
	{Ext: ".lasso", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Lasso
	{Ext: ".mako", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Mako
	{Ext: ".gms", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // GAMS
	{Ext: ".icl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Clean
	{Ext: ".arc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Arc
	{Ext: ".wast", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // WebAssembly
	{Ext: ".spin", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Propeller Spin
	{Ext: ".po", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Gettext Catalog
	{Ext: ".rsc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Rascal
	{Ext: ".x10", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // X10
	{Ext: ".ston", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // STON
	{Ext: ".muf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // MUF
	{Ext: ".dats", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // ATS
	{Ext: ".adb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Ada
	{Ext: ".nc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // nesC
	{Ext: ".rhtml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // RHTML
	{Ext: ".nu", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Nu
	{Ext: ".flf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // FIGlet Font
	{Ext: ".asp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // ASP
	{Ext: ".nl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // NL
	{Ext: ".nsi", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // NSIS
	{Ext: ".vala", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Vala
	{Ext: ".ecl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // ECL
	{Ext: ".bsv", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Bluespec
	{Ext: ".axs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // NetLinx
	{Ext: ".6pl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Perl 6
	{Ext: ".qml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // QML
	{Ext: ".eml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // EML
	{Ext: ".sls", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // SaltStack
	{Ext: ".brd", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // KiCad Legacy Layout
	{Ext: ".fish", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // fish
	{Ext: ".fan", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Fantom
	{Ext: ".pike", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Pike
	{Ext: ".s", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // Unix Assembly
	{Ext: ".xc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // XC
	{Ext: ".ijs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // J
	{Ext: ".asciidoc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},      // AsciiDoc
	{Ext: ".for", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Formatted
	{Ext: ".tex", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // TeX
	{Ext: ".pep", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Pep8
	{Ext: ".tla", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // TLA
	{Ext: ".r", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // R
	{Ext: ".lua", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Lua
	{Ext: ".xs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // XS
	{Ext: ".smali", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Smali
	{Ext: ".bal", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Ballerina
	{Ext: ".upc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Unified Parallel C
	{Ext: ".md", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Markdown
	{Ext: ".ps", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // PostScript
	{Ext: ".tea", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Tea
	{Ext: ".sql", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // SQLPL
	{Ext: ".feature", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},       // Gherkin
	{Ext: ".styl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Stylus
	{Ext: ".wisp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // wisp
	{Ext: ".gdb", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // GDB
	{Ext: ".apib", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // API Blueprint
	{Ext: ".as", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // ActionScript
	{Ext: ".diff", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Diff
	{Ext: ".cppobjdump", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},    // Cpp-ObjDump
	{Ext: ".twig", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Twig
	{Ext: ".zep", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Zephir
	{Ext: ".click", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // Click
	{Ext: ".obj", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Wavefront Object
	{Ext: ".dm", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // DM
	{Ext: ".ik", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Ioke
	{Ext: ".gp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Gnuplot
	{Ext: ".jsonld", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // JSONLD
	{Ext: ".dwl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // DataWeave
	{Ext: ".ecl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // ECLiPSe
	{Ext: ".p", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // OpenEdge ABL
	{Ext: ".hx", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Haxe
	{Ext: ".sfd", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Spline Font Database
	{Ext: ".mu", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // mupad
	{Ext: ".soy", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Closure Templates
	{Ext: ".pan", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Pan
	{Ext: ".lookml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // LookML
	{Ext: ".mod", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Linux Kernel Module
	{Ext: ".txl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // TXL
	{Ext: ".liquid", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Liquid
	{Ext: ".nim", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Nim
	{Ext: ".dockerfile", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},    // Dockerfile
	{Ext: ".maxpat", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Max
	{Ext: ".lisp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Common Lisp
	{Ext: ".kit", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Kit
	{Ext: ".nix", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Nix
	{Ext: ".sss", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // SugarSS
	{Ext: ".toml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // TOML
	{Ext: ".xquery", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // XQuery
	{Ext: ".nit", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Nit
	{Ext: ".pov", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // POV-Ray SDL
	{Ext: ".ll", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // LLVM
	{Ext: ".E", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},             // E
	{Ext: ".parrot", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // Parrot
	{Ext: ".gf", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Grammatical Framework
	{Ext: ".asc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // AGS Script
	{Ext: ".mumps", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // M
	{Ext: ".psc", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Papyrus
	{Ext: ".cpp", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // C++
	{Ext: ".rnh", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // RUNOFF
	{Ext: ".mss", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // CartoCSS
	{Ext: ".cwl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Common Workflow Language
	{Ext: ".shader", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // ShaderLab
	{Ext: ".pkl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Pickle
	{Ext: ".sco", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Csound Score
	{Ext: ".rbbas", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // REALbasic
	{Ext: ".ejs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // EJS
	{Ext: ".moon", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // MoonScript
	{Ext: ".pwn", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Pawn
	{Ext: ".jisonlex", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},      // Jison Lex
	{Ext: ".aug", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Augeas
	{Ext: ".slim", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},          // Slim
	{Ext: ".irclog", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},        // IRC log
	{Ext: ".fs", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},            // Filterscript
	{Ext: ".bro", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Bro
	{Ext: ".omgrofl", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},       // Omgrofl
	{Ext: ".gml", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Game Maker Language
	{Ext: ".rkt", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Racket
	{Ext: ".nlogo", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},         // NetLogo
	{Ext: ".pod", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},           // Pod
	{Ext: ".litcoffee", FileType: FileTypeLanguage, Supported: true, PreviewType: FileTypeText},     // Literate CoffeeScript
}
