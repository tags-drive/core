package web

type icon struct {
	// Name of an icon without ".svg"
	iconFilename string

	extensions []string
	filenames  []string
}

// icons is a list of supported file icons (from https://github.com/PKief/vscode-material-icon-theme)
var icons = [...]icon{
	{
		iconFilename: "html",
		extensions:   []string{"html", "htm", "xhtml", "html_vm", "asp"},
	},
	{
		iconFilename: "pug",
		extensions:   []string{"jade", "pug"},
	},
	{
		iconFilename: "markdown",
		extensions:   []string{"md", "markdown", "rst"},
	},
	{
		iconFilename: "blink",
		extensions:   []string{"blink"},
	},
	{
		iconFilename: "css",
		extensions:   []string{"css"},
	},
	{
		iconFilename: "sass",
		extensions:   []string{"scss", "sass"},
	},
	{
		iconFilename: "less",
		extensions:   []string{"less"},
	},
	{
		iconFilename: "json",
		extensions:   []string{"json", "tsbuildinfo"},
		filenames: []string{".jscsrc",
			".jshintrc",
			"tsconfig.json",
			"tslint.json",
			"composer.lock",
			".jsbeautifyrc",
			".esformatter",
			"cdp.pid",
		},
	},
	{
		iconFilename: "jinja",
		extensions:   []string{"jinja", "jinja2", "j2"},
	},
	{
		iconFilename: "sublime",
		extensions:   []string{"sublime-project", "sublime-workspace"},
	},
	{
		iconFilename: "yaml",
		extensions:   []string{"yaml", "YAML-tmLanguage", "yml"},
	},
	{
		iconFilename: "xml",
		extensions:   []string{"xml", "plist", "xsd", "dtd", "xsl", "xslt", "resx", "iml", "xquery", "tmLanguage", "manifest", "project"},
		filenames:    []string{".htaccess"},
	},
	{
		iconFilename: "image",
		extensions:   []string{"png", "jpeg", "jpg", "gif", "ico", "tif", "tiff", "psd", "psb", "ami", "apx", "bmp", "bpg", "brk", "cur", "dds", "dng", "exr", "fpx", "gbr", "img", "jbig2", "jb2", "jng", "jxr", "pbm", "pgf", "pic", "raw", "webp", "eps"},
	},
	{
		iconFilename: "javascript",
		extensions:   []string{"js", "esx", "mjs"},
	},
	{
		iconFilename: "react",
		extensions:   []string{"jsx"},
	},
	{
		iconFilename: "react_ts",
		extensions:   []string{"tsx"},
	},
	{
		iconFilename: "routing",
		extensions:   []string{"routing.ts", "routing.tsx", "routing.js", "routing.jsx", "routes.ts", "routes.tsx", "routes.js", "routes.jsx"},
		filenames:    []string{"router.js", "router.jsx", "router.ts", "router.tsx", "routes.js", "routes.jsx", "routes.ts", "routes.tsx"},
	},
	{
		iconFilename: "settings",
		extensions: []string{"ini", "dlc", "dll", "config", "conf", "properties", "prop", "settings", "option", "props", "toml", "prefs", "sln.dotsettings",
			"sln.dotsettings.user",
			"cfg"},
		filenames: []string{".jshintignore",
			".buildignore",
			".mrconfig",
			".yardopts",
			"manifest.mf",
			".clang-format",
			".clang-tidy",
		},
	},
	{
		iconFilename: "typescript",
		extensions:   []string{"ts"},
	},
	{
		iconFilename: "typescript-def",
		extensions:   []string{"d.ts"},
	},
	{
		iconFilename: "markojs",
		extensions:   []string{"marko"},
	},
	{
		iconFilename: "pdf",
		extensions:   []string{"pdf"},
	},
	{
		iconFilename: "table",
		extensions:   []string{"xlsx", "xls", "csv", "tsv"},
	},
	{
		iconFilename: "vscode",
		extensions:   []string{"vscodeignore", "vsixmanifest", "vsix", "code-workplace"},
	},
	{
		iconFilename: "visualstudio",
		extensions: []string{"csproj", "ruleset", "sln", "suo", "vb", "vbs", "vcxitems", "vcxitems.filters",
			"vcxproj", "vcxproj.filters",
		},
	},
	{
		iconFilename: "database",
		extensions:   []string{"pdb", "sql", "pks", "pkb", "accdb", "mdb", "sqlite", "pgsql", "postgres", "psql"},
	},
	{
		iconFilename: "csharp",
		extensions:   []string{"cs", "csx"},
	},
	{
		iconFilename: "zip",
		extensions:   []string{"zip", "tar", "gz", "xz", "br", "bzip2", "gzip", "brotli", "7z", "rar", "tgz"},
	},
	{
		iconFilename: "exe",
		extensions:   []string{"exe", "msi"},
	},
	{
		iconFilename: "java",
		extensions:   []string{"java", "jar", "jsp"},
	},
	{
		iconFilename: "c",
		extensions:   []string{"c", "m", "i", "mi"},
	},
	{
		iconFilename: "h",
		extensions:   []string{"h"},
	},
	{
		iconFilename: "cpp",
		extensions:   []string{"cc", "cpp", "cxx", "c++", "cp", "mm", "mii", "ii"},
	},
	{
		iconFilename: "hpp",
		extensions:   []string{"hh", "hpp", "hxx", "h++", "hp", "tcc", "inl"},
	},
	{
		iconFilename: "go",
		extensions:   []string{"go"},
	},
	{
		iconFilename: "go-mod",
		filenames:    []string{"go.mod", "go.sum"},
	},
	{
		iconFilename: "python",
		extensions:   []string{"py"},
	},
	{
		iconFilename: "python-misc",
		extensions:   []string{"pyc", "whl"},
		filenames:    []string{"requirements.txt", "pipfile", ".python-version", "manifest.in"},
	},
	{
		iconFilename: "url",
		extensions:   []string{"url"},
	},
	{
		iconFilename: "console",
		extensions:   []string{"sh", "ksh", "csh", "tcsh", "zsh", "bash", "bat", "cmd", "awk", "fish"},
	},
	{
		iconFilename: "powershell",
		extensions:   []string{"ps1", "psm1", "psd1", "ps1xml", "psc1", "pssc"},
	},
	{
		iconFilename: "gradle",
		extensions:   []string{"gradle"},
		filenames:    []string{"gradle.properties", "gradlew", "gradle-wrapper.properties"},
	},
	{
		iconFilename: "word",
		extensions:   []string{"doc", "docx", "rtf"},
	},
	{
		iconFilename: "certificate",
		extensions:   []string{"cer", "cert", "crt"},
		filenames: []string{"license", "license.md",
			"license.txt",
			"licence", "licence.md",
			"licence.txt",
			"unlicense", "unlicense.md",
			"unlicense.txt",
		},
	},
	{
		iconFilename: "key",
		extensions:   []string{"pub", "key", "pem", "asc", "gpg"},
		filenames:    []string{".htpasswd"},
	},
	{
		iconFilename: "font",
		extensions:   []string{"woff", "woff2", "ttf", "eot", "suit", "otf", "bmap", "fnt", "odttf", "ttc", "font", "fonts", "sui", "ntf", "mrf"},
	},
	{
		iconFilename: "lib",
		extensions:   []string{"lib", "bib"},
	},
	{
		iconFilename: "ruby",
		extensions:   []string{"rb", "erb"},
	},
	{
		iconFilename: "gemfile",
		filenames:    []string{"gemfile"},
	},
	{
		iconFilename: "fsharp",
		extensions:   []string{"fs", "fsx", "fsi", "fsproj"},
	},
	{
		iconFilename: "swift",
		extensions:   []string{"swift"},
	},
	{
		iconFilename: "arduino",
		extensions:   []string{"ino"},
	},
	{
		iconFilename: "docker",
		extensions:   []string{"dockerignore", "dockerfile"},
		filenames:    []string{"dockerfile", "docker-compose.yml", "docker-compose.yaml", "docker-compose.dev.yml", "docker-compose.local.yml", "docker-compose.ci.yml", "docker-compose.override.yml", "docker-compose.staging.yml", "docker-compose.prod.yml", "docker-compose.production.yml", "docker-compose.test.yml"},
	},
	{
		iconFilename: "tex",
		extensions:   []string{"tex", "cls", "sty", "dtx", "ltx"},
	},
	{
		iconFilename: "powerpoint",
		extensions:   []string{"pptx", "ppt", "pptm", "potx", "potm", "ppsx", "ppsm", "pps", "ppam", "ppa"},
	},
	{
		iconFilename: "video",
		extensions:   []string{"webm", "mkv", "flv", "vob", "ogv", "ogg", "gifv", "avi", "mov", "qt", "wmv", "yuv", "rm", "rmvb", "mp4", "m4v", "mpg", "mp2", "mpeg", "mpe", "mpv", "m2v"},
	},
	{
		iconFilename: "virtual",
		extensions:   []string{"vdi", "vbox", "vbox-prev"},
	},
	{
		iconFilename: "email",
		extensions:   []string{"ics"},
		filenames:    []string{".mailmap"},
	},
	{
		iconFilename: "audio",
		extensions:   []string{"mp3", "flac", "m4a", "wma", "aiff"},
	},
	{
		iconFilename: "coffee",
		extensions:   []string{"coffee", "cson", "iced"},
	},
	{
		iconFilename: "document",
		extensions:   []string{"txt"},
	},
	{
		iconFilename: "graphql",
		extensions:   []string{"graphql", "gql"},
		filenames:    []string{".graphqlconfig"},
	},
	{
		iconFilename: "rust",
		extensions:   []string{"rs"},
	},
	{
		iconFilename: "raml",
		extensions:   []string{"raml"},
	},
	{
		iconFilename: "xaml",
		extensions:   []string{"xaml"},
	},
	{
		iconFilename: "haskell",
		extensions:   []string{"hs"},
	},
	{
		iconFilename: "kotlin",
		extensions:   []string{"kt", "kts"},
	},
	{
		iconFilename: "git",
		extensions:   []string{"patch"},
		filenames:    []string{".gitignore", ".gitconfig", ".gitattributes", ".gitmodules", ".gitkeep", "git-history"},
	},
	{
		iconFilename: "lua",
		extensions:   []string{"lua"},
		filenames:    []string{".luacheckrc"},
	},
	{
		iconFilename: "clojure",
		extensions:   []string{"clj", "cljs", "cljc"},
	},
	{
		iconFilename: "groovy",
		extensions:   []string{"groovy"},
	},
	{
		iconFilename: "r",
		extensions:   []string{"r", "rmd"},
		filenames:    []string{".Rhistory"},
	},
	{
		iconFilename: "dart",
		extensions:   []string{"dart"},
	},
	{
		iconFilename: "actionscript",
		extensions:   []string{"as"},
	},
	{
		iconFilename: "mxml",
		extensions:   []string{"mxml"},
	},
	{
		iconFilename: "autohotkey",
		extensions:   []string{"ahk"},
	},
	{
		iconFilename: "flash",
		extensions:   []string{"swf"},
	},
	{
		iconFilename: "swc",
		extensions:   []string{"swc"},
	},
	{
		iconFilename: "cmake",
		extensions:   []string{"cmake"},
		filenames:    []string{"cmakelists.txt", "cmakecache.txt"},
	},
	{
		iconFilename: "assembly",
		extensions:   []string{"asm", "a51", "inc", "nasm", "s", "ms", "agc", "ags", "aea", "argus", "mitigus", "binsource"},
	},
	{
		iconFilename: "vue",
		extensions:   []string{"vue"},
	},
	{
		iconFilename: "vue-config",
		filenames:    []string{"vue.config.js", "vue.config.ts"},
	},
	{
		iconFilename: "ocaml",
		extensions:   []string{"ml", "mli", "cmx"},
	},
	{
		iconFilename: "javascript-map",
		extensions:   []string{"js.map", "mjs.map"},
	},
	{
		iconFilename: "css-map",
		extensions:   []string{"css.map"},
	},
	{
		iconFilename: "lock",
		extensions:   []string{"lock"},
	},
	{
		iconFilename: "handlebars",
		extensions:   []string{"hbs", "mustache"},
	},
	{
		iconFilename: "perl",
		extensions:   []string{"pl", "pm"},
	},
	{
		iconFilename: "haxe",
		extensions:   []string{"hx"},
	},
	{
		iconFilename: "test-ts",
		extensions:   []string{"spec.ts", "e2e-spec.ts", "test.ts", "ts.snap"},
	},
	{
		iconFilename: "test-jsx",
		extensions:   []string{"spec.tsx", "test.tsx", "tsx.snap", "spec.jsx", "test.jsx", "jsx.snap"},
	},
	{
		iconFilename: "test-js",
		extensions:   []string{"spec.js", "e2e-spec.js", "test.js", "js.snap"},
	},
	{
		iconFilename: "puppet",
		extensions:   []string{"pp"},
	},
	{
		iconFilename: "elixir",
		extensions:   []string{"ex", "exs", "eex", "leex"},
	},
	{
		iconFilename: "livescript",
		extensions:   []string{"ls"},
	},
	{
		iconFilename: "erlang",
		extensions:   []string{"erl"},
	},
	{
		iconFilename: "twig",
		extensions:   []string{"twig"},
	},
	{
		iconFilename: "julia",
		extensions:   []string{"jl"},
	},
	{
		iconFilename: "elm",
		extensions:   []string{"elm"},
	},
	{
		iconFilename: "purescript",
		extensions:   []string{"pure", "purs"},
	},
	{
		iconFilename: "smarty",
		extensions:   []string{"tpl"},
	},
	{
		iconFilename: "stylus",
		extensions:   []string{"styl"},
	},
	{
		iconFilename: "reason",
		extensions:   []string{"re", "rei"},
	},
	{
		iconFilename: "bucklescript",
		extensions:   []string{"cmj"},
	},
	{
		iconFilename: "merlin",
		extensions:   []string{"merlin"},
	},
	{
		iconFilename: "verilog",
		extensions:   []string{"v", "vhd", "sv", "svh"},
	},
	{
		iconFilename: "mathematica",
		extensions:   []string{"nb"},
	},
	{
		iconFilename: "wolframlanguage",
		extensions:   []string{"wl", "wls"},
	},
	{
		iconFilename: "nunjucks",
		extensions:   []string{"njk", "nunjucks"},
	},
	{
		iconFilename: "robot",
		extensions:   []string{"robot"},
	},
	{
		iconFilename: "solidity",
		extensions:   []string{"sol"},
	},
	{
		iconFilename: "autoit",
		extensions:   []string{"au3"},
	},
	{
		iconFilename: "haml",
		extensions:   []string{"haml"},
	},
	{
		iconFilename: "yang",
		extensions:   []string{"yang"},
	},
	{
		iconFilename: "mjml",
		extensions:   []string{"mjml"},
	},
	{
		iconFilename: "terraform",
		extensions:   []string{"tf", "tf.json", "tfvars", "tfstate"},
	},
	{
		iconFilename: "laravel",
		extensions:   []string{"blade.php", "inky.php"},
	},
	{
		iconFilename: "applescript",
		extensions:   []string{"applescript"},
	},
	{
		iconFilename: "cake",
		extensions:   []string{"cake"},
	},
	{
		iconFilename: "cucumber",
		extensions:   []string{"feature"},
	},
	{
		iconFilename: "nim",
		extensions:   []string{"nim", "nimble"},
	},
	{
		iconFilename: "apiblueprint",
		extensions:   []string{"apib", "apiblueprint"},
	},
	{
		iconFilename: "riot",
		extensions:   []string{"riot", "tag"},
	},
	{
		iconFilename: "vfl",
		extensions:   []string{"vfl"},
		filenames:    []string{".vfl"},
	},
	{
		iconFilename: "kl",
		extensions:   []string{"kl"},
		filenames:    []string{".kl"},
	},
	{
		iconFilename: "postcss",
		extensions:   []string{"pcss", "sss"},
		filenames:    []string{"postcss.config.js", ".postcssrc.js", ".postcssrc", ".postcssrc.json", ".postcssrc.yml"},
	},
	{
		iconFilename: "todo",
		extensions:   []string{"todo"},
	},
	{
		iconFilename: "coldfusion",
		extensions:   []string{"cfml", "cfc", "lucee", "cfm"},
	},
	{
		iconFilename: "cabal",
		extensions:   []string{"cabal"},
	},
	{
		iconFilename: "nix",
		extensions:   []string{"nix"},
	},
	{
		iconFilename: "slim",
		extensions:   []string{"slim"},
	},
	{
		iconFilename: "http",
		extensions:   []string{"http", "rest"},
	},
	{
		iconFilename: "restql",
		extensions:   []string{"rql", "restql"},
	},
	{
		iconFilename: "kivy",
		extensions:   []string{"kv"},
	},
	{
		iconFilename: "graphcool",
		extensions:   []string{"graphcool"},
		filenames:    []string{"project.graphcool"},
	},
	{
		iconFilename: "sbt",
		extensions:   []string{"sbt"},
	},
	{
		iconFilename: "webpack",
		filenames:    []string{"webpack.js", "webpack.ts", "webpack.base.js", "webpack.base.ts", "webpack.config.js", "webpack.config.ts", "webpack.common.js", "webpack.common.ts", "webpack.config.common.js", "webpack.config.common.ts", "webpack.config.common.babel.js", "webpack.config.common.babel.ts", "webpack.dev.js", "webpack.dev.ts", "webpack.config.dev.js", "webpack.config.dev.ts", "webpack.config.dev.babel.js", "webpack.config.dev.babel.ts", "webpack.prod.js", "webpack.prod.ts", "webpack.server.js", "webpack.server.ts", "webpack.client.js", "webpack.client.ts", "webpack.config.server.js", "webpack.config.server.ts", "webpack.config.client.js", "webpack.config.client.ts", "webpack.config.production.babel.js", "webpack.config.production.babel.ts", "webpack.config.prod.babel.js", "webpack.config.prod.babel.ts", "webpack.config.prod.js", "webpack.config.prod.ts", "webpack.config.production.js", "webpack.config.production.ts", "webpack.config.staging.js", "webpack.config.staging.ts", "webpack.config.babel.js", "webpack.config.babel.ts", "webpack.config.base.babel.js", "webpack.config.base.babel.ts", "webpack.config.base.js", "webpack.config.base.ts", "webpack.config.staging.babel.js", "webpack.config.staging.babel.ts", "webpack.config.coffee", "webpack.config.test.js", "webpack.config.test.ts", "webpack.config.vendor.js", "webpack.config.vendor.ts", "webpack.config.vendor.production.js", "webpack.config.vendor.production.ts", "webpack.test.js", "webpack.test.ts", "webpack.dist.js", "webpack.dist.ts", "webpackfile.js", "webpackfile.ts"},
	},
	{
		iconFilename: "ionic",
		filenames:    []string{"ionic.config.json", ".io-config.json"},
	},
	{
		iconFilename: "gulp",
		filenames:    []string{"gulpfile.js", "gulpfile.ts", "gulpfile.babel.js"},
	},
	{
		iconFilename: "nodejs",
		filenames:    []string{"package.json", "package-lock.json", ".nvmrc", ".esmrc"},
	},
	{
		iconFilename: "npm",
		filenames:    []string{".npmignore", ".npmrc"},
	},
	{
		iconFilename: "yarn",
		filenames:    []string{".yarnrc", "yarn.lock", ".yarnclean", ".yarn-integrity", "yarn-error.log"},
	},
	{
		iconFilename: "android",
		filenames:    []string{"androidmanifest.xml"},
	},
	{
		iconFilename: "tune",
		extensions:   []string{"env"},
		filenames:    []string{".env.example", ".env.local", ".env.dev", ".env.development", ".env.prod", ".env.production", ".env.staging", ".env.preview", ".env.test", ".env.development.local", ".env.production.local", ".env.test.local"},
	},
	{
		iconFilename: "babel",
		filenames:    []string{".babelrc", ".babelrc.js", "babel.config.js"},
	},
	{
		iconFilename: "contributing",
		filenames:    []string{"contributing.md"},
	},
	{
		iconFilename: "readme",
		filenames:    []string{"readme.md", "readme.txt", "readme"},
	},
	{
		iconFilename: "changelog",
		filenames:    []string{"changelog", "changelog.md", "changelog.txt"},
	},
	{
		iconFilename: "credits",
		filenames:    []string{"credits", "credits.txt", "credits.md"},
	},
	{
		iconFilename: "authors",
		filenames:    []string{"authors", "authors.md", "authors.txt"},
	},
	{
		iconFilename: "flow",
		filenames:    []string{".flowconfig"},
	},
	{
		iconFilename: "favicon",
		filenames:    []string{"favicon.ico"},
	},
	{
		iconFilename: "karma",
		filenames:    []string{"karma.conf.js", "karma.conf.ts", "karma.conf.coffee", "karma.config.js", "karma.config.ts", "karma-main.js", "karma-main.ts"},
	},
	{
		iconFilename: "bithound",
		filenames:    []string{".bithoundrc"},
	},
	{
		iconFilename: "appveyor",
		filenames:    []string{".appveyor.yml", "appveyor.yml"},
	},
	{
		iconFilename: "travis",
		filenames:    []string{".travis.yml"},
	},
	{
		iconFilename: "protractor",
		filenames:    []string{"protractor.conf.js", "protractor.conf.ts", "protractor.conf.coffee", "protractor.config.js", "protractor.config.ts"},
	},
	{
		iconFilename: "fusebox",
		filenames:    []string{"fuse.js"},
	},
	{
		iconFilename: "heroku",
		filenames:    []string{"procfile", "procfile.windows"},
	},
	{
		iconFilename: "editorconfig",
		filenames:    []string{".editorconfig"},
	},
	{
		iconFilename: "gitlab",
		extensions:   []string{"gitlab-ci.yml"},
	},
	{
		iconFilename: "bower",
		filenames:    []string{".bowerrc", "bower.json"},
	},
	{
		iconFilename: "eslint",
		filenames:    []string{".eslintrc.js", ".eslintrc.yaml", ".eslintrc.yml", ".eslintrc.json", ".eslintrc", ".eslintignore"},
	},
	{
		iconFilename: "conduct",
		filenames:    []string{"code_of_conduct.md", "code_of_conduct.txt"},
	},
	{
		iconFilename: "watchman",
		filenames:    []string{".watchmanconfig"},
	},
	{
		iconFilename: "aurelia",
		filenames:    []string{"aurelia.json"},
	},
	{
		iconFilename: "mocha",
		filenames:    []string{"mocha.opts", ".mocharc.yml", ".mocharc.yaml", ".mocharc.js", ".mocharc.json", ".mocharc.jsonc"},
	},
	{
		iconFilename: "jenkins",
		filenames:    []string{"jenkinsfile"},
		extensions:   []string{"jenkinsfile", "jenkins"},
	},
	{
		iconFilename: "firebase",
		filenames:    []string{"firebase.json", ".firebaserc"},
	},
	{
		iconFilename: "rollup",
		filenames:    []string{"rollup.config.js", "rollup.config.ts", "rollup-config.js", "rollup-config.ts", "rollup.config.common.js", "rollup.config.common.ts", "rollup.config.base.js", "rollup.config.base.ts", "rollup.config.prod.js", "rollup.config.prod.ts", "rollup.config.dev.js", "rollup.config.dev.ts", "rollup.config.prod.vendor.js", "rollup.config.prod.vendor.ts"},
	},
	{
		iconFilename: "hack",
		filenames:    []string{".hhconfig"},
	},
	{
		iconFilename: "stylelint",
		filenames:    []string{".stylelintrc", "stylelint.config.js", ".stylelintrc.json", ".stylelintrc.yaml", ".stylelintrc.yml", ".stylelintrc.js", ".stylelintignore"},
	},
	{
		iconFilename: "code-climate",
		filenames:    []string{".codeclimate.yml"},
	},
	{
		iconFilename: "prettier",
		filenames:    []string{".prettierrc", "prettier.config.js", ".prettierrc.js", ".prettierrc.json", ".prettierrc.yaml", ".prettierrc.yml", ".prettierignore"},
	},
	{
		iconFilename: "nodemon",
		filenames:    []string{"nodemon.json", "nodemon-debug.json"},
	},
	{
		iconFilename: "webhint",
		filenames:    []string{".hintrc"},
	},
	{
		iconFilename: "browserlist",
		filenames:    []string{"browserslist", ".browserslistrc"},
	},
	{
		iconFilename: "crystal",
		extensions:   []string{"cr", "ecr"},
	},
	{
		iconFilename: "snyk",
		filenames:    []string{".snyk"},
	},
	{
		iconFilename: "drone",
		extensions:   []string{"drone.yml"},
		filenames:    []string{".drone.yml"},
	},
	{
		iconFilename: "cuda",
		extensions:   []string{"cu", "cuh"},
	},
	{
		iconFilename: "log",
		extensions:   []string{"log"},
	},
	{
		iconFilename: "dotjs",
		extensions:   []string{"def", "dot", "jst"},
	},
	{
		iconFilename: "ejs",
		extensions:   []string{"ejs"},
	},
	{
		iconFilename: "sequelize",
		filenames:    []string{".sequelizerc"},
	},
	{
		iconFilename: "gatsby",
		filenames:    []string{"gatsby.config.js", "gatsby-config.js", "gatsby-node.js", "gatsby-browser.js", "gatsby-ssr.js"},
	},
	{
		iconFilename: "wakatime",
		filenames:    []string{".wakatime-project"},
		extensions:   []string{".wakatime-project"},
	},
	{
		iconFilename: "circleci",
		filenames:    []string{"circle.yml"},
	},
	{
		iconFilename: "cloudfoundry",
		filenames:    []string{".cfignore"},
	},
	{
		iconFilename: "grunt",
		filenames:    []string{"gruntfile.js", "gruntfile.ts", "gruntfile.coffee", "gruntfile.babel.js", "gruntfile.babel.ts", "gruntfile.babel.coffee"},
	},
	{
		iconFilename: "jest",
		filenames:    []string{"jest.config.js", "jest.config.ts", "jest.config.json", "jest.setup.js", "jest.setup.ts", "jest.json", ".jestrc", "jest.teardown.js"},
	},
	{
		iconFilename: "processing",
		extensions:   []string{"pde"},
	},
	{
		iconFilename: "storybook",
		extensions:   []string{"stories.js", "stories.jsx", "story.js", "story.jsx", "stories.ts", "stories.tsx", "story.ts", "story.tsx"},
	},
	{
		iconFilename: "wepy",
		extensions:   []string{"wpy"},
	},
	{
		iconFilename: "fastlane",
		filenames:    []string{"fastfile", "appfile"},
	},
	{
		iconFilename: "hcl",
		extensions:   []string{"hcl"},
	},
	{
		iconFilename: "helm",
		filenames:    []string{".helmignore"},
	},
	{
		iconFilename: "san",
		extensions:   []string{"san"},
	},
	{
		iconFilename: "wallaby",
		filenames:    []string{"wallaby.js", "wallaby.conf.js"},
	},
	{
		iconFilename: "django",
		extensions:   []string{"djt"},
	},
	{
		iconFilename: "stencil",
		filenames:    []string{"stencil.config.js", "stencil.config.ts"},
	},
	{
		iconFilename: "red",
		extensions:   []string{"red"},
	},
	{
		iconFilename: "makefile",
		filenames:    []string{"makefile"},
	},
	{
		iconFilename: "foxpro",
		extensions:   []string{"fxp", "prg"},
	},
	{
		iconFilename: "i18n",
		extensions:   []string{"pot", "po", "mo"},
	},
	{
		iconFilename: "webassembly",
		extensions:   []string{"wat", "wasm"},
	},
	{
		iconFilename: "semantic-release",
		filenames:    []string{".releaserc", "release.config.js"},
	},
	{
		iconFilename: "bitbucket",
		filenames:    []string{"bitbucket-pipelines.yaml", "bitbucket-pipelines.yml"},
	},
	{
		iconFilename: "jupyter",
		extensions:   []string{"ipynb"},
	},
	{
		iconFilename: "d",
		extensions:   []string{"d"},
	},
	{
		iconFilename: "mdx",
		extensions:   []string{"mdx"},
	},
	{
		iconFilename: "ballerina",
		extensions:   []string{"bal", "balx"},
	},
	{
		iconFilename: "racket",
		extensions:   []string{"rkt"},
	},
	{
		iconFilename: "bazel",
		extensions:   []string{"bzl", "bazel"},
		filenames:    []string{".bazelignore", ".bazelrc"},
	},
	{
		iconFilename: "mint",
		extensions:   []string{"mint"},
	},
	{
		iconFilename: "velocity",
		extensions:   []string{"vm", "fhtml", "vtl"},
	},
	{
		iconFilename: "godot",
		extensions:   []string{"gd"},
	},
	{
		iconFilename: "godot-assets",
		extensions:   []string{"godot", "tres", "tscn"},
	},
	{
		iconFilename: "azure-pipelines",
		filenames:    []string{"azure-pipelines.yml"},
	},
	{
		iconFilename: "azure",
		extensions:   []string{"azcli"},
	},
	{
		iconFilename: "vagrant",
		filenames:    []string{"vagrantfile"},
		extensions:   []string{"vagrantfile"},
	},
	{
		iconFilename: "prisma",
		filenames:    []string{"prisma.yml"},
		extensions:   []string{"prisma"},
	},
	{
		iconFilename: "razor",
		extensions:   []string{"cshtml", "vbhtml"},
	},
	{
		iconFilename: "asciidoc",
		extensions:   []string{"ad", "adoc", "asciidoc"},
	},
	{
		iconFilename: "istanbul",
		filenames:    []string{".nycrc", ".nycrc.json"},
	},
	{
		iconFilename: "edge",
		extensions:   []string{"edge"},
	},
	{
		iconFilename: "scheme",
		extensions:   []string{"ss", "scm"},
	},
	{
		iconFilename: "tailwindcss",
		filenames:    []string{"tailwind.js", "tailwind.config.js"},
	},
	{
		iconFilename: "3d",
		extensions:   []string{"stl", "obj", "ac"},
	},
	{
		iconFilename: "buildkite",
		filenames:    []string{"buildkite.yml", "buildkite.yaml"},
	},
	{
		iconFilename: "netlify",
		filenames:    []string{"netlify.toml"},
	},
	{
		iconFilename: "svg",
		extensions:   []string{"svg"},
	},
	{
		iconFilename: "svelte",
		extensions:   []string{"svelte"},
	},
	{
		iconFilename: "vim",
		extensions:   []string{"vimrc", "gvimrc", "exrc"},
	},
	{
		iconFilename: "nest",
		filenames:    []string{"nest-cli.json", ".nest-cli.json", "nestconfig.json", ".nestconfig.json"},
	},
	{
		iconFilename: "moonscript",
		extensions:   []string{"moon"},
	},
	{
		iconFilename: "percy",
		filenames:    []string{".percy.yml"},
	},
	{
		iconFilename: "gitpod",
		filenames:    []string{".gitpod.yml"},
	},
}
