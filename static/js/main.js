// For every fetch request we have to use `credentials: "same-origin"` to pass Cookie

const sortType = {
    name: "name",
    size: "size",
    time: "time"
};

const sortOrder = {
    asc: "asc",
    desc: "desc"
};

// store contains global data
var store = {
    state: {
        msg: "Test",
        allFiles: [],
        allTags: [],
        opacity: 1
    },
    updateFiles: function() {
        fetch("/api/files", {
            method: "GET",
            credentials: "same-origin"
        })
            .then(data => data.json())
            .then(files => this.setFiles(files))
            .catch(err => console.log(err));
    },
    updateTags: function() {
        fetch("/api/tags", {
            method: "GET",
            credentials: "same-origin"
        })
            .then(data => data.json())
            .then(tags => {
                this.state.allTags = tags;
            })
            .catch(err => console.log(err));
    },
    setFiles: function(files) {
        // Change time from "2018-08-23T22:48:59.0459184+03:00" to "23-08-2018 22:48"
        for (let i in files) {
            files[i].addTime = new Date(files[i].addTime).format(
                "dd-mm-yyyy HH:MM"
            );
        }
        this.state.allFiles = files;
    }
};

// Init should be called onload
function Init() {
    store.updateTags();
    store.updateFiles();
    recentFiles.update();
}

function updateStore() {
    store.updateFiles();
    store.updateTags();
}

// Upload block
var uploader = new Vue({
    el: "#upload-block",
    data: {
        sharedState: store.state,
        counter: 0 // for definition did user drag file into div. If counter > 0, user dragged file.
    },
    created() {
        // Add listeners
        document.ondragenter = () => {
            this.counter++;
        };
        document.ondragleave = () => {
            this.counter--;
        };
        document.ondrop = () => {
            this.counter = 0;
        };
        setInterval(() => {
            if (this.counter == 0) {
                this.sharedState.opacity = 1;
            } else {
                this.sharedState.opacity = 0.3;
            }
        }, 20);
    },
    methods: {
        upload: function(event) {
            var formData = new FormData();

            for (file of event.dataTransfer.files) {
                formData.append("files", file, file.name);
            }

            fetch("/api/files", {
                body: formData,
                method: "POST",
                credentials: "same-origin"
            })
                .then(res => res.json())
                .then(log => {
                    console.log(log);
                    this.logs = log;
                    // Update list of files
                    updateStore();
                })
                .catch(err => console.log(err));
        }
    }
});

// Main block
var mainBlock = new Vue({
    el: "#main-block",
    data: {
        sharedState: store.state,
        opacity: 1,

        sortByNameModeAsc: true,
        sortBySizeModeAsc: true,
        sortByTimeModeAsc: true,
        lastSortType: sortType.name
    },
    methods: {
        // Context menu
        showContextMenu: function(event, file) {
            contextMenu.setFile(file);
            contextMenu.showMenu(event.x, event.y);
        },
        // Sort
        sortByName: function() {
            if (this.lastSortType == sortType.name) {
                this.sortByNameModeAsc = !this.sortByNameModeAsc;
            } else {
                // Use default settings
                this.resetSortTypes();
            }
            this.lastSortType = sortType.name;

            let type = sortType.name,
                order = this.sortByNameModeAsc ? sortOrder.asc : sortOrder.desc;

            searchBar.advancedSearch(type, order);
        },
        sortBySize: function() {
            if (this.lastSortType == sortType.size) {
                this.sortBySizeModeAsc = !this.sortBySizeModeAsc;
            } else {
                // Use default settings
                this.resetSortTypes();
            }
            this.lastSortType = sortType.size;

            let type = sortType.size,
                order = this.sortBySizeModeAsc ? sortOrder.asc : sortOrder.desc;

            searchBar.advancedSearch(type, order);
        },
        sortByTime: function() {
            if (this.lastSortType == sortType.time) {
                this.sortByTimeModeAsc = !this.sortByTimeModeAsc;
            } else {
                // Use default settings
                this.resetSortTypes();
            }
            this.lastSortType = sortType.time;

            let type = sortType.time,
                order = this.sortByTimeModeAsc ? sortOrder.asc : sortOrder.desc;

            searchBar.advancedSearch(type, order);
        },
        resetSortTypes: function() {
            this.sortByNameModeAsc = true;
            this.sortBySizeModeAsc = true;
            this.sortByTimeModeAsc = true;
        }
    },
    template: `
	<table :style="{'opacity': sharedState.opacity}" class="file-table" style="width:100%;">
			<tr style="position: sticky; top: 100px;">
				<th></th>
				<th>
					Filename
					<i class="material-icons" id="sortByNameIcon" @click="sortByName" :style="[sortByNameModeAsc ? {'transform': 'scale(1, 1)'} : {'transform': 'scale(1, -1)'}]" style="font-size: 20px; cursor: pointer;">
						sort
					</i>
				</th>
				<th>Tags</th>
				<th>
					Size (MB)
					<i class="material-icons" id="sortByNameSize" @click="sortBySize" :style="[sortBySizeModeAsc ? {'transform': 'scale(1, 1)'} : {'transform': 'scale(1, -1)'}]" style="transform: scale(1, 1); font-size: 20px; cursor: pointer;">
						sort
					</i>
				</th>
				<th>
					Time of adding
					<i class="material-icons" id="sortByNameTime" @click="sortByTime" :style="[sortByTimeModeAsc ? {'transform': 'scale(1, 1)'} : {'transform': 'scale(1, -1)'}]" style="transform: scale(1, 1); font-size: 20px; cursor: pointer;">
						sort
					</i>
				</th>
			</tr>
			<files v-for="file in sharedState.allFiles" :file="file"></files>
	</table>`
});

var contextMenu = new Vue({
    el: "#context-menu",
    mixins: [VueClickaway.mixin], // from vue-clickaway
    data: {
        file: null,
        // Style
        top: "0px",
        left: "0px",
        show: false,
        // File changing
        newName: "",
        newTags: [],
        description: "",
        // For calculation of position
        divWidth: 140,
        divHeight: 125
    },
    methods: {
        setFile: function(file) {
            this.file = file;
        },
        showMenu: function(x, y) {
            const offset = 10;
            x += offset;
            y += offset;
            if (x + this.divWidth > window.innerWidth) {
                x -= offset * 2;
                x -= this.divWidth;
            }
            if (y + this.divHeight > window.innerHeight) {
                y -= offset * 2;
                y -= this.divHeight;
            }
            this.left = x + "px";
            this.top = y + "px";
            this.show = true;
        },
        hideMenu: function() {
            this.show = false;
        },
        changeName: function() {
            this.show = false;
            console.log("Change name");
        },
        changeTags: function() {
            this.show = false;
            console.log("Change tags");
        },
        changeDescription: function() {
            this.show = false;
            console.log("Change description");
        },
        deleteFile: function() {
            this.show = false;
            console.log("Delete");
        }
    }
});

// Search bar
var searchBar = new Vue({
    el: "#top-bar__search",
    data: {
        sharedState: store.state,
        tagForAdding: "",
        pickedTags: [],
        text: "",
        selectedMode: "And"
    },
    methods: {
        search: function() {
            let params = new URLSearchParams();
            // tags
            if (this.pickedTags.length != 0) {
                let tags = [];
                for (let tag of this.pickedTags) {
                    tags.push(tag.name);
                }
                params.append("tags", tags.join(","));
            }
            // search
            if (this.text != "") {
                params.append("search", this.text);
            }
            // sort
            params.append("sort", sortType.name);
            // order
            params.append("order", sortOrder.asc);
            // mode
            params.append("mode", this.selectedMode.toLowerCase());

            fetch("/api/files?" + params, {
                method: "GET",
                credentials: "same-origin"
            })
                .then(data => data.json())
                .then(files => {
                    store.setFiles(files);
                    // Reset sortParams
                    mainBlock.resetSortTypes();
                });
        },
        advancedSearch: function(sType, sOrder) {
            let params = new URLSearchParams();
            // tags
            if (this.pickedTags.length != 0) {
                let tags = [];
                for (let tag of this.pickedTags) {
                    tags.push(tag.name);
                }
                params.append("tags", tags.join(","));
            }
            // search
            if (this.text != "") {
                params.append("search", this.text);
            }
            // sort
            params.append("sort", sType);
            // order
            params.append("order", sOrder);
            // mode
            params.append("mode", this.selectedMode.toLowerCase());

            fetch("/api/files?" + params, {
                method: "GET",
                credentials: "same-origin"
            })
                .then(data => data.json())
                .then(files => store.setFiles(files));
        },
        deleteTagFromSearch: function(name) {
            let index = -1;
            for (i in this.pickedTags) {
                if (this.pickedTags[i].name == name) {
                    index = i;
                    break;
                }
            }
            if (index == -1) {
                return;
            }

            // Remove an element
            this.pickedTags.splice(index, 1);
        },
        addTag: function() {
            // Check is there the tag
            for (let tag of this.sharedState.allTags) {
                if (tag.name == this.tagForAdding) {
                    let alreadyHas = false;
                    // Check was tag already picked
                    for (let tag of this.pickedTags) {
                        if (tag.name == this.tagForAdding) {
                            alreadyHas = true;
                            break;
                        }
                    }
                    if (!alreadyHas) {
                        this.tagForAdding = "";
                        this.pickedTags.push(tag);
                    }

                    break;
                }
            }
        }
    }
});

// Recent files
var recentFiles = new Vue({
    el: "#recentFiles",
    data: {
        recentFiles: []
    },
    methods: {
        update: function() {
            fetch("/api/files/recent", {
                method: "GET",
                credentials: "same-origin"
            })
                .then(data => data.json())
                .then(files => (this.recentFiles = files));
        }
    }
});
