package extensions

// List with all supported types of files
var extensionsList = []Ext{
	// Archives
	{Ext: ".7z", FileType: FileTypeArchive, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".pkg", FileType: FileTypeArchive, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".rar", FileType: FileTypeArchive, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".zip", FileType: FileTypeArchive, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".tar.gz", FileType: FileTypeArchive, Supported: false, PreviewType: PreviewTypeUnsupported},

	// Audio
	{Ext: ".aif", FileType: FileTypeAudio, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".mpa", FileType: FileTypeAudio, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".wma", FileType: FileTypeAudio, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".wpl", FileType: FileTypeAudio, Supported: false, PreviewType: PreviewTypeUnsupported},
	// supproted
	{Ext: ".ogg", FileType: FileTypeAudio, Supported: true, PreviewType: PreviewTypeAudioOGG},
	{Ext: ".wav", FileType: FileTypeAudio, Supported: true, PreviewType: PreviewTypeAudioWAV},
	{Ext: ".mp3", FileType: FileTypeAudio, Supported: true, PreviewType: PreviewTypeAudioMP3},

	// Image
	{Ext: ".bmp", FileType: FileTypeImage, Supported: true, PreviewType: PreviewTypeImage},
	{Ext: ".gif", FileType: FileTypeImage, Supported: true, PreviewType: PreviewTypeImage},
	{Ext: ".ico", FileType: FileTypeImage, Supported: true, PreviewType: PreviewTypeImage},
	{Ext: ".jpg", FileType: FileTypeImage, Supported: true, PreviewType: PreviewTypeImage},
	{Ext: ".png", FileType: FileTypeImage, Supported: true, PreviewType: PreviewTypeImage},
	{Ext: ".svg", FileType: FileTypeImage, Supported: true, PreviewType: PreviewTypeImage},
	{Ext: ".jpeg", FileType: FileTypeImage, Supported: true, PreviewType: PreviewTypeImage},

	// Video
	{Ext: ".avi", FileType: FileTypeVideo, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".mkv", FileType: FileTypeVideo, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".mov", FileType: FileTypeVideo, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".mpg", FileType: FileTypeVideo, Supported: false, PreviewType: PreviewTypeUnsupported},
	{Ext: ".mpeg", FileType: FileTypeVideo, Supported: false, PreviewType: PreviewTypeUnsupported},
	// supported
	{Ext: ".mp4", FileType: FileTypeVideo, Supported: true, PreviewType: PreviewTypeVideoMP4},
	{Ext: ".webm", FileType: FileTypeVideo, Supported: true, PreviewType: PreviewTypeVideoWebM},

	// Text (from https://github.com/github/linguist/blob/master/lib/linguist/languages.yml)
	{Ext: ".cfg", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // HAProxy
	{Ext: ".m", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Mercury
	{Ext: ".pic", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Pic
	{Ext: ".pb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // PureBasic
	{Ext: ".mm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Objective-C++
	{Ext: ".asn", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // ASN.1
	{Ext: ".d", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // DTrace
	{Ext: ".self", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Self
	{Ext: ".marko", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Marko
	{Ext: ".lean", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Lean
	{Ext: ".nl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // NewLisp
	{Ext: ".conllu", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // CoNLL-U
	{Ext: ".reb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Rebol
	{Ext: ".robot", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // RobotFramework
	{Ext: ".xpm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // X PixMap
	{Ext: ".properties", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},    // Java Properties
	{Ext: ".pas", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Pascal
	{Ext: ".brs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Brightscript
	{Ext: ".owl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Web Ontology Language
	{Ext: ".q", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // HiveQL
	{Ext: ".cshtml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // HTML+Razor
	{Ext: ".e", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Eiffel
	{Ext: ".befunge", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},       // Befunge
	{Ext: ".raml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // RAML
	{Ext: ".opa", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Opa
	{Ext: ".ex", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Elixir
	{Ext: ".ur", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // UrWeb
	{Ext: ".mq4", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // MQL4
	{Ext: ".mq5", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // MQL5
	{Ext: ".agda", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Agda
	{Ext: ".thrift", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Thrift
	{Ext: ".xm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Logos
	{Ext: ".srt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // SRecode Template
	{Ext: ".jinja", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // HTML+Django
	{Ext: ".sass", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Sass
	{Ext: ".pbt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // PowerBuilder
	{Ext: ".edn", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // edn
	{Ext: ".sublime-build", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText}, // JSON with Comments
	{Ext: ".blade", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Blade
	{Ext: ".pp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Puppet
	{Ext: ".cr", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Crystal
	{Ext: ".m", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Objective-C
	{Ext: ".bison", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Bison
	{Ext: ".druby", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Mirah
	{Ext: ".gn", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // GN
	{Ext: ".ebuild", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Gentoo Ebuild
	{Ext: ".j", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Objective-J
	{Ext: ".bdf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Glyph Bitmap Distribution Format
	{Ext: ".ftl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // FreeMarker
	{Ext: ".f90", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Fortran
	{Ext: ".swift", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Swift
	{Ext: ".phtml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // HTML+PHP
	{Ext: ".ipynb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Jupyter Notebook
	{Ext: ".tpl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Smarty
	{Ext: ".cp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Component Pascal
	{Ext: ".cob", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // COBOL
	{Ext: ".go", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Go
	{Ext: ".ring", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Ring
	{Ext: ".php", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // PHP
	{Ext: ".csd", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Csound Document
	{Ext: ".bsl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // 1C Enterprise
	{Ext: ".i3", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Modula-3
	{Ext: ".http", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // HTTP
	{Ext: ".ms", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // MAXScript
	{Ext: ".raw", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Raw token data
	{Ext: ".mak", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Makefile
	{Ext: ".bb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // BitBake
	{Ext: ".ml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // OCaml
	{Ext: ".textile", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},       // Textile
	{Ext: ".m4", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // M4
	{Ext: ".sed", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // sed
	{Ext: ".v", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Verilog
	{Ext: ".do", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Stata
	{Ext: ".gitignore", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // Ignore List
	{Ext: ".red", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Red
	{Ext: ".ahk", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // AutoHotkey
	{Ext: ".yar", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // YARA
	{Ext: ".d-objdump", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // D-ObjDump
	{Ext: ".xpl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // XProc
	{Ext: ".xsp-config", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},    // XPages
	{Ext: ".volt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Volt
	{Ext: ".as", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // AngelScript
	{Ext: ".lsl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // LSL
	{Ext: ".jade", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Pug
	{Ext: ".rs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // RenderScript
	{Ext: ".minid", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // MiniD
	{Ext: ".t", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Terra
	{Ext: ".coffee", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // CoffeeScript
	{Ext: ".html", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // HTML
	{Ext: ".awk", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Awk
	{Ext: ".dot", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Graphviz (DOT)
	{Ext: ".org", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Org
	{Ext: ".idr", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Idris
	{Ext: ".pcss", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // PostCSS
	{Ext: ".pgsql", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // PLpgSQL
	{Ext: ".agc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Apollo Guidance Computer
	{Ext: ".mod", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Modula-2
	{Ext: ".c", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // C
	{Ext: ".myt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Myghty
	{Ext: ".hlsl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // HLSL
	{Ext: ".p4", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // P4
	{Ext: ".thy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Isabelle
	{Ext: ".eex", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // HTML+EEX
	{Ext: ".scaml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Scaml
	{Ext: ".clw", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Clarion
	{Ext: ".sh-session", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},    // ShellSession
	{Ext: ".pod", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Pod 6
	{Ext: ".graphql", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},       // GraphQL
	{Ext: ".lhs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Literate Haskell
	{Ext: ".erb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // HTML+ERB
	{Ext: ".au3", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // AutoIt
	{Ext: ".gs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Genie
	{Ext: ".moo", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Moocode
	{Ext: ".hcl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // HCL
	{Ext: ".chs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // C2hs Haskell
	{Ext: ".jsp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Java Server Pages
	{Ext: ".sp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // SourcePawn
	{Ext: ".cmake", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // CMake
	{Ext: ".clj", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Clojure
	{Ext: ".shen", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Shen
	{Ext: ".wdl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // wdl
	{Ext: ".hxml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // HXML
	{Ext: ".yml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // YAML
	{Ext: ".tcl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Tcl
	{Ext: ".nf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Nextflow
	{Ext: ".ini", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // INI
	{Ext: ".xtend", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Xtend
	{Ext: ".pytb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Python traceback
	{Ext: ".epj", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Ecere Projects
	{Ext: ".glsl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // GLSL
	{Ext: ".uc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // UnrealScript
	{Ext: ".eb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Easybuild
	{Ext: ".jq", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // JSONiq
	{Ext: ".vhdl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // VHDL
	{Ext: ".pyx", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Cython
	{Ext: ".sv", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // SystemVerilog
	{Ext: ".objdump", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},       // ObjDump
	{Ext: ".smt2", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // SMT
	{Ext: ".regexp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Regular Expression
	{Ext: ".gitconfig", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // Git Config
	{Ext: ".json5", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // JSON5
	{Ext: ".yasnippet", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // YASnippet
	{Ext: ".anim", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Unity3D Asset
	{Ext: ".b", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Limbo
	{Ext: ".rexx", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // REXX
	{Ext: ".lfe", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // LFE
	{Ext: ".oxygene", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},       // Oxygene
	{Ext: ".bb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // BlitzBasic
	{Ext: ".asm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Assembly
	{Ext: ".flex", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // JFlex
	{Ext: ".cs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // C#
	{Ext: ".t", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Turing
	{Ext: ".pde", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Processing
	{Ext: ".fx", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // FLUX
	{Ext: ".webidl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // WebIDL
	{Ext: ".json", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // JSON
	{Ext: ".asc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Public Key
	{Ext: ".matlab", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // MATLAB
	{Ext: ".ls", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // LiveScript
	{Ext: ".stan", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Stan
	{Ext: ".mtml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // MTML
	{Ext: ".pl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Prolog
	{Ext: ".jsx", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // JSX
	{Ext: ".proto", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Protocol Buffer
	{Ext: ".ch", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Charity
	{Ext: ".xml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // XML
	{Ext: ".svg", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // SVG
	{Ext: ".mathematica", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},   // Mathematica
	{Ext: ".lol", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // LOLCODE
	{Ext: ".ooc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // ooc
	{Ext: ".xslt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // XSLT
	{Ext: ".rmd", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // RMarkdown
	{Ext: ".afm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Adobe Font Metrics
	{Ext: ".mask", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Mask
	{Ext: ".sch", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // KiCad Schematic
	{Ext: ".em", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // EmberScript
	{Ext: ".lvproj", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // LabVIEW
	{Ext: ".n", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Nemerle
	{Ext: ".cu", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Cuda
	{Ext: ".krl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // KRL
	{Ext: ".vim", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Vim script
	{Ext: ".pony", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Pony
	{Ext: ".sci", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Scilab
	{Ext: ".1", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Roff Manpage
	{Ext: ".rpy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Ren'Py
	{Ext: ".sparql", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // SPARQL
	{Ext: ".applescript", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},   // AppleScript
	{Ext: ".txt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Text
	{Ext: ".sage", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Sage
	{Ext: ".ck", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // ChucK
	{Ext: ".g4", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // ANTLR
	{Ext: ".fs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // F#
	{Ext: ".ls", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // LoomScript
	{Ext: ".fy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Fancy
	{Ext: ".fst", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // F*
	{Ext: ".pir", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Parrot Internal Representation
	{Ext: ".st", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Smalltalk
	{Ext: ".ice", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Slice
	{Ext: ".monkey", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Monkey
	{Ext: ".pogo", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // PogoScript
	{Ext: ".el", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Emacs Lisp
	{Ext: ".js", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // JavaScript
	{Ext: ".pro", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // QMake
	{Ext: ".rs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Rust
	{Ext: ".abap", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // ABAP
	{Ext: ".pasm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Parrot Assembly
	{Ext: ".cw", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Redcode
	{Ext: ".sl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Slash
	{Ext: ".l", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // PicoLisp
	{Ext: ".spec", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // RPM Spec
	{Ext: ".erl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Erlang
	{Ext: ".mms", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Module Management System
	{Ext: ".dae", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // COLLADA
	{Ext: ".scm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Scheme
	{Ext: ".nut", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Squirrel
	{Ext: ".py", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Python
	{Ext: ".nanorc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // nanorc
	{Ext: ".latte", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Latte
	{Ext: ".ne", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Nearley
	{Ext: ".iss", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Inno Setup
	{Ext: ".ebnf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // EBNF
	{Ext: ".ipf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // IGOR Pro
	{Ext: ".chpl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Chapel
	{Ext: ".coq", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Coq
	{Ext: ".dylan", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Dylan
	{Ext: ".lagda", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Literate Agda
	{Ext: ".sch", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Eagle
	{Ext: ".gradle", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Gradle
	{Ext: ".clp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // CLIPS
	{Ext: ".axs.erb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},       // NetLinx+ERB
	{Ext: ".eclass", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Gentoo Eclass
	{Ext: ".xbm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // X BitMap
	{Ext: ".als", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Alloy
	{Ext: ".groovy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Groovy
	{Ext: ".w", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // CWeb
	{Ext: ".ol", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Jolie
	{Ext: ".pls", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // PLSQL
	{Ext: ".purs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // PureScript
	{Ext: ".jl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Julia
	{Ext: ".bf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // HyPhy
	{Ext: ".q", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // q
	{Ext: ".hs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Haskell
	{Ext: ".ncl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // NCL
	{Ext: ".vb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Visual Basic
	{Ext: ".io", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Io
	{Ext: ".rg", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Rouge
	{Ext: ".haml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Haml
	{Ext: ".djs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Dogescript
	{Ext: ".ps1", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // PowerShell
	{Ext: ".ts", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // TypeScript
	{Ext: ".dart", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Dart
	{Ext: ".edc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Edje Data Collection
	{Ext: ".vcl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // VCL
	{Ext: ".zig", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Zig
	{Ext: ".ceylon", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Ceylon
	{Ext: ".fr", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Frege
	{Ext: ".pro", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // IDL
	{Ext: ".g", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // GAP
	{Ext: ".aj", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // AspectJ
	{Ext: ".sh", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Shell
	{Ext: ".orc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Csound
	{Ext: ".tcsh", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Tcsh
	{Ext: ".prg", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // xBase
	{Ext: ".elm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Elm
	{Ext: ".jison", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Jison
	{Ext: ".x", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // RPC
	{Ext: ".desktop", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},       // desktop
	{Ext: ".sc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // SuperCollider
	{Ext: ".nginxconf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // Nginx
	{Ext: ".re", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Reason
	{Ext: ".yang", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // YANG
	{Ext: ".com", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // DIGITAL Command Language
	{Ext: ".sas", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // SAS
	{Ext: ".ninja", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Ninja
	{Ext: ".grace", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Grace
	{Ext: ".cl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // OpenCL
	{Ext: ".d", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // D
	{Ext: ".creole", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Creole
	{Ext: ".kt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Kotlin
	{Ext: ".opal", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Opal
	{Ext: ".8xp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // TI Program
	{Ext: ".ML", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Standard ML
	{Ext: ".cfc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // ColdFusion CFC
	{Ext: ".bat", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Batchfile
	{Ext: ".oz", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Oz
	{Ext: ".ox", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Ox
	{Ext: ".gsp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Groovy Server Pages
	{Ext: ".roff", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Roff
	{Ext: ".rl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Ragel
	{Ext: ".gs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Gosu
	{Ext: ".handlebars", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},    // Handlebars
	{Ext: ".less", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Less
	{Ext: ".zone", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // DNS Zone
	{Ext: ".pd", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Pure Data
	{Ext: ".ecr", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // HTML+ECR
	{Ext: ".kicad_pcb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // KiCad Layout
	{Ext: ".ld", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Linker Script
	{Ext: ".b", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Brainfuck
	{Ext: ".f", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Filebench WML
	{Ext: ".apl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // APL
	{Ext: ".hh", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Hack
	{Ext: ".toc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // World of Warcraft Addon Data
	{Ext: ".numpy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // NumPy
	{Ext: ".sqf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // SQF
	{Ext: ".glf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Glyph
	{Ext: ".fea", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // OpenType Feature File
	{Ext: ".cy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Cycript
	{Ext: ".java", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Java
	{Ext: ".scala", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Scala
	{Ext: ".scad", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // OpenSCAD
	{Ext: ".apacheconf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},    // ApacheConf
	{Ext: ".pl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Perl
	{Ext: ".asy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // LTspice Symbol
	{Ext: ".mediawiki", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // MediaWiki
	{Ext: ".vue", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Vue
	{Ext: ".gd", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // GDScript
	{Ext: ".gbr", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Gerber Image
	{Ext: ".capnp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Cap'n Proto
	{Ext: ".factor", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Factor
	{Ext: ".reg", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Windows Registry Entries
	{Ext: ".darcspatch", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},    // Darcs Patch
	{Ext: ".fth", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Forth
	{Ext: ".asy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Asymptote
	{Ext: ".hy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Hy
	{Ext: ".j", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Jasmin
	{Ext: ".ec", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // eC
	{Ext: ".scss", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // SCSS
	{Ext: ".cls", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Apex
	{Ext: ".l", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Lex
	{Ext: ".rb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Ruby
	{Ext: ".ly", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // LilyPond
	{Ext: ".cl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Cool
	{Ext: ".zimpl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Zimpl
	{Ext: ".kid", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Genshi
	{Ext: ".golo", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Golo
	{Ext: ".cson", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // CSON
	{Ext: ".sql", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // SQL
	{Ext: ".metal", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Metal
	{Ext: ".gml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Graph Modeling Language
	{Ext: ".md", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // GCC Machine Description
	{Ext: ".ni", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Inform 7
	{Ext: ".lgt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Logtalk
	{Ext: ".mo", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Modelica
	{Ext: ".m4", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // M4Sugar
	{Ext: ".boo", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Boo
	{Ext: ".csv", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // CSV
	{Ext: ".eq", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // EQ
	{Ext: ".mtl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Wavefront Material
	{Ext: ".css", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // CSS
	{Ext: ".uno", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Uno
	{Ext: ".ttl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Turtle
	{Ext: ".c-objdump", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // C-ObjDump
	{Ext: ".rdoc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // RDoc
	{Ext: ".abnf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // ABNF
	{Ext: ".ampl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // AMPL
	{Ext: ".cfm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // ColdFusion
	{Ext: ".cirru", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Cirru
	{Ext: ".rst", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // reStructuredText
	{Ext: ".hb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Harbour
	{Ext: ".y", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Yacc
	{Ext: ".g", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // G-code
	{Ext: ".xojo_code", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // Xojo
	{Ext: ".srt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // SubRip Text
	{Ext: ".bmx", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // BlitzMax
	{Ext: ".pig", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // PigLatin
	{Ext: ".tl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Type Language
	{Ext: ".lasso", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Lasso
	{Ext: ".mako", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Mako
	{Ext: ".gms", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // GAMS
	{Ext: ".icl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Clean
	{Ext: ".arc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Arc
	{Ext: ".wast", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // WebAssembly
	{Ext: ".spin", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Propeller Spin
	{Ext: ".po", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Gettext Catalog
	{Ext: ".rsc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Rascal
	{Ext: ".x10", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // X10
	{Ext: ".ston", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // STON
	{Ext: ".muf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // MUF
	{Ext: ".dats", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // ATS
	{Ext: ".adb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Ada
	{Ext: ".nc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // nesC
	{Ext: ".rhtml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // RHTML
	{Ext: ".nu", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Nu
	{Ext: ".flf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // FIGlet Font
	{Ext: ".asp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // ASP
	{Ext: ".nl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // NL
	{Ext: ".nsi", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // NSIS
	{Ext: ".vala", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Vala
	{Ext: ".ecl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // ECL
	{Ext: ".bsv", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Bluespec
	{Ext: ".axs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // NetLinx
	{Ext: ".6pl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Perl 6
	{Ext: ".qml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // QML
	{Ext: ".eml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // EML
	{Ext: ".sls", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // SaltStack
	{Ext: ".brd", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // KiCad Legacy Layout
	{Ext: ".fish", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // fish
	{Ext: ".fan", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Fantom
	{Ext: ".pike", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Pike
	{Ext: ".s", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // Unix Assembly
	{Ext: ".xc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // XC
	{Ext: ".ijs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // J
	{Ext: ".asciidoc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},      // AsciiDoc
	{Ext: ".for", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Formatted
	{Ext: ".tex", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // TeX
	{Ext: ".pep", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Pep8
	{Ext: ".tla", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // TLA
	{Ext: ".r", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // R
	{Ext: ".lua", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Lua
	{Ext: ".xs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // XS
	{Ext: ".smali", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Smali
	{Ext: ".bal", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Ballerina
	{Ext: ".upc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Unified Parallel C
	{Ext: ".md", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Markdown
	{Ext: ".ps", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // PostScript
	{Ext: ".tea", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Tea
	{Ext: ".sql", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // SQLPL
	{Ext: ".feature", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},       // Gherkin
	{Ext: ".styl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Stylus
	{Ext: ".wisp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // wisp
	{Ext: ".gdb", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // GDB
	{Ext: ".apib", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // API Blueprint
	{Ext: ".as", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // ActionScript
	{Ext: ".diff", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Diff
	{Ext: ".cppobjdump", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},    // Cpp-ObjDump
	{Ext: ".twig", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Twig
	{Ext: ".zep", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Zephir
	{Ext: ".click", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // Click
	{Ext: ".obj", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Wavefront Object
	{Ext: ".dm", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // DM
	{Ext: ".ik", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Ioke
	{Ext: ".gp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Gnuplot
	{Ext: ".jsonld", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // JSONLD
	{Ext: ".dwl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // DataWeave
	{Ext: ".ecl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // ECLiPSe
	{Ext: ".p", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // OpenEdge ABL
	{Ext: ".hx", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Haxe
	{Ext: ".sfd", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Spline Font Database
	{Ext: ".mu", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // mupad
	{Ext: ".soy", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Closure Templates
	{Ext: ".pan", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Pan
	{Ext: ".lookml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // LookML
	{Ext: ".mod", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Linux Kernel Module
	{Ext: ".txl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // TXL
	{Ext: ".liquid", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Liquid
	{Ext: ".nim", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Nim
	{Ext: ".dockerfile", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},    // Dockerfile
	{Ext: ".maxpat", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Max
	{Ext: ".lisp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Common Lisp
	{Ext: ".kit", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Kit
	{Ext: ".nix", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Nix
	{Ext: ".sss", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // SugarSS
	{Ext: ".toml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // TOML
	{Ext: ".xquery", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // XQuery
	{Ext: ".nit", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Nit
	{Ext: ".pov", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // POV-Ray SDL
	{Ext: ".ll", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // LLVM
	{Ext: ".E", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},             // E
	{Ext: ".parrot", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // Parrot
	{Ext: ".gf", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Grammatical Framework
	{Ext: ".asc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // AGS Script
	{Ext: ".mumps", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // M
	{Ext: ".psc", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Papyrus
	{Ext: ".cpp", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // C++
	{Ext: ".rnh", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // RUNOFF
	{Ext: ".mss", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // CartoCSS
	{Ext: ".cwl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Common Workflow Language
	{Ext: ".shader", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // ShaderLab
	{Ext: ".pkl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Pickle
	{Ext: ".sco", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Csound Score
	{Ext: ".rbbas", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // REALbasic
	{Ext: ".ejs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // EJS
	{Ext: ".moon", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // MoonScript
	{Ext: ".pwn", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Pawn
	{Ext: ".jisonlex", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},      // Jison Lex
	{Ext: ".aug", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Augeas
	{Ext: ".slim", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},          // Slim
	{Ext: ".irclog", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},        // IRC log
	{Ext: ".fs", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},            // Filterscript
	{Ext: ".bro", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Bro
	{Ext: ".omgrofl", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},       // Omgrofl
	{Ext: ".gml", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Game Maker Language
	{Ext: ".rkt", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Racket
	{Ext: ".nlogo", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},         // NetLogo
	{Ext: ".pod", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},           // Pod
	{Ext: ".litcoffee", FileType: FileTypeLanguage, Supported: true, PreviewType: PreviewTypeText},     // Literate CoffeeScript
}
