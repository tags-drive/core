// For every fetch request we have to use `credentials: "same-origin"` to pass Cookie

/* Global */
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
        allFiles: [],
        allTags: {},
        opacity: 1,
        showDropLayer: true // when we show modal-window with tags showDropLayer is false
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
            files[i].addTime = new Date(files[i].addTime).format("dd-mm-yyyy HH:MM");
        }
        this.state.allFiles = files;
    }
};

// Init should be called onload
function Init() {
    store.updateTags();
    store.updateFiles();
    leftBar.update();
}

function updateStore() {
    store.updateFiles();
    store.updateTags();
}

/* Main instances */

// Top bar
var topBar = new Vue({
    el: "#top-bar",
    data: {
        sharedState: store.state,
        tagForAdding: "",
        pickedTags: [],
        text: "",
        selectedMode: "And"
    },
    methods: {
        search: function() {
            return {
                usual: () => {
                    let params = new URLSearchParams();
                    // tags
                    if (this.pickedTags.length != 0) {
                        let tags = [];
                        for (let tag of this.pickedTags) {
                            tags.push(tag.id);
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
                            mainBlock.sort().reset();
                        });
                },
                advanced: (sType, sOrder) => {
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
                }
            };
        },
        input: function() {
            return {
                tags: {
                    add: () => {
                        // TODO
                        // Check is there the tag
                        for (let id in this.sharedState.allTags) {
                            if (this.sharedState.allTags[id].name == this.tagForAdding) {
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
                                    this.pickedTags.push(this.sharedState.allTags[id]);
                                }

                                break;
                            }
                        }
                    },
                    delete: tagID => {
                        let index = -1;
                        for (i in this.pickedTags) {
                            if (this.pickedTags[i].id == tagID) {
                                index = i;
                                break;
                            }
                        }
                        if (index == -1) {
                            return;
                        }

                        // Remove an element
                        this.pickedTags.splice(index, 1);
                    }
                }
            };
        },
        settings: function() {
            return {
                tags: () => modalWindow.showWindow().globalTags()
            };
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
        // Sorts
        sort: function() {
            return {
                byName: () => {
                    if (this.lastSortType == sortType.name) {
                        this.sortByNameModeAsc = !this.sortByNameModeAsc;
                    } else {
                        // Use default settings
                        this.sort().reset();
                    }
                    this.lastSortType = sortType.name;

                    let type = sortType.name,
                        order = this.sortByNameModeAsc ? sortOrder.asc : sortOrder.desc;

                    topBar.search().advanced(type, order);
                },
                bySize: () => {
                    if (this.lastSortType == sortType.size) {
                        this.sortBySizeModeAsc = !this.sortBySizeModeAsc;
                    } else {
                        // Use default settings
                        this.sort().reset();
                    }
                    this.lastSortType = sortType.size;

                    let type = sortType.size,
                        order = this.sortBySizeModeAsc ? sortOrder.asc : sortOrder.desc;

                    topBar.search().advanced(type, order);
                },
                byTime: () => {
                    if (this.lastSortType == sortType.time) {
                        this.sortByTimeModeAsc = !this.sortByTimeModeAsc;
                    } else {
                        // Use default settings
                        this.sort().reset();
                    }
                    this.lastSortType = sortType.time;

                    let type = sortType.time,
                        order = this.sortByTimeModeAsc ? sortOrder.asc : sortOrder.desc;

                    topBar.search().advanced(type, order);
                },
                reset: () => {
                    this.sortByNameModeAsc = true;
                    this.sortBySizeModeAsc = true;
                    this.sortByTimeModeAsc = true;
                }
            };
        }
    },
    template: `
	<table :style="{'opacity': sharedState.opacity}" class="file-table" style="width:100%;">
		<tr style="position: sticky; top: 100px;">
			<th></th>
			<th>
				Filename
				<i class="material-icons" id="sortByNameIcon" @click="sort().byName()" :style="[sortByNameModeAsc ? {'transform': 'scale(1, 1)'} : {'transform': 'scale(1, -1)'}]" style="font-size: 20px; cursor: pointer;">
					sort
				</i>
			</th>
			<th>Tags</th>
			<th>
				Size (MB)
				<i class="material-icons" id="sortByNameSize" @click="sort().bySize()" :style="[sortBySizeModeAsc ? {'transform': 'scale(1, 1)'} : {'transform': 'scale(1, -1)'}]" style="transform: scale(1, 1); font-size: 20px; cursor: pointer;">
					sort
				</i>
			</th>
			<th>
				Time of adding
				<i class="material-icons" id="sortByNameTime" @click="sort().byTime()" :style="[sortByTimeModeAsc ? {'transform': 'scale(1, 1)'} : {'transform': 'scale(1, -1)'}]" style="transform: scale(1, 1); font-size: 20px; cursor: pointer;">
					sort
				</i>
			</th>
		</tr>
		<files v-for="file in sharedState.allFiles" :file="file" :allTags="store.state.allTags"></files>
	</table>`
});

// Left bar
var leftBar = new Vue({
    el: "#left-bar",
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

/* Secondary instances */

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
            if (store.state.showDropLayer) {
                this.counter++;
            }
        };
        document.ondragleave = () => {
            if (store.state.showDropLayer) {
                this.counter--;
            }
        };
        document.ondrop = () => {
            if (store.state.showDropLayer) {
                this.counter = 0;
            }
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

// Context menu (right click on a file)
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
        // UI
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
        // Options of context menu
        changeName: function() {
            this.show = false;
            modalWindow.showWindow().renaming(this.file);
        },
        changeTags: function() {
            this.show = false;
            modalWindow.showWindow().fileTags(this.file);
        },
        changeDescription: function() {
            this.show = false;
            modalWindow.showWindow().description(this.file);
        },
        deleteFile: function() {
            this.show = false;
            modalWindow.showWindow().deleting(this.file);
        }
    }
});

// Modal window
// It's called from a context menu
var modalWindow = new Vue({
    el: "#modal-window",
    data: {
        file: null,
        show: false,
        error: "",
        sharedState: store.state,
        // Modes
        renameMode: false,
        tagsMode: false,
        descriptionMode: false,
        deleteMode: false,
        globalTagsMode: false,
        // For files API
        fileNewData: {
            newFilename: "",
            unusedTags: [],
            newTags: [],
            newDescription: ""
        },
        // For tags API
        newTag: {}
    },
    methods: {
        // UI
        showWindow: function() {
            return {
                renaming: file => {
                    this.file = file;
                    this.renameMode = true;
                    this.fileNewData.newFilename = file.filename;

                    this.show = true;
                },
                fileTags: file => {
                    store.state.showDropLayer = false;

                    this.fileNewData.newTags = [];
                    this.fileNewData.unusedTags = [];

                    for (let id in this.sharedState.allTags) {
                        if (file.tags.includes(Number(id))) {
                            this.fileNewData.newTags.push(this.sharedState.allTags[id]);
                        } else {
                            this.fileNewData.unusedTags.push(this.sharedState.allTags[id]);
                        }
                    }

                    this.file = file;
                    this.tagsMode = true;

                    this.show = true;
                },
                description: file => {
                    this.file = file;
                    this.fileNewData.unusedTags;
                    this.descriptionMode = true;

                    this.show = true;
                },
                deleting: file => {
                    this.file = file;
                    this.deleteMode = true;

                    this.show = true;
                },
                globalTags: () => {
                    this.globalTagsMode = true;

                    this.show = true;
                }
            };
        },
        hideWindow: function() {
            this.renameMode = false;
            this.tagsMode = false;
            this.descriptionMode = false;
            this.deleteMode = false;
            this.show = false;
            store.state.showDropLayer = true;
        },
        // Drag and drop
        tagsDragAndDrop: function() {
            return {
                addToFile: ev => {
                    let tagID = Number(ev.dataTransfer.getData("tagName"));
                    let index = -1;
                    for (i in this.fileNewData.unusedTags) {
                        if (this.fileNewData.unusedTags[i].id == tagID) {
                            index = i;
                            break;
                        }
                    }
                    if (index == -1) {
                        return;
                    }
                    this.fileNewData.newTags.push(this.fileNewData.unusedTags[index]);
                    this.fileNewData.unusedTags.splice(index, 1);
                },
                delFromFile: ev => {
                    let tagID = ev.dataTransfer.getData("tagName");
                    let index = -1;
                    for (i in this.fileNewData.newTags) {
                        if (this.fileNewData.newTags[i].id == tagID) {
                            index = i;
                            break;
                        }
                    }
                    if (index == -1) {
                        return;
                    }
                    this.fileNewData.unusedTags.push(this.fileNewData.newTags[index]);
                    this.fileNewData.newTags.splice(index, 1);
                }
            };
        },
        // Files API
        filesAPI: function() {
            return {
                rename: () => {
                    let params = new URLSearchParams();
                    params.append("file", this.file.filename);
                    params.append("new-name", this.fileNewData.newFilename);

                    fetch("/api/files", {
                        method: "PUT",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (resp.status >= 400 && resp.status < 600) {
                                // TODO: return resp.text(). How to do?
                                throw new Error("TODO");
                            }
                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => {
                            this.error = err;
                            console.log(err);
                        });
                },
                updateTags: () => {
                    let params = new URLSearchParams();
                    let tags = this.fileNewData.newTags.map(tag => tag.id);
                    params.append("file", this.file.filename);
                    params.append("tags", tags.join(","));

                    fetch("/api/files", {
                        method: "PUT",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (resp.status >= 400 && resp.status < 600) {
                                // TODO: return resp.text(). How to do?
                                throw new Error("TODO");
                            }
                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => {
                            this.error = err;
                            console.log(err);
                        });
                },
                updateDescription: () => {
                    let params = new URLSearchParams();
                    params.append("file", this.file.filename);
                    params.append("description", this.fileNewData.newDescription);

                    fetch("/api/files", {
                        method: "PUT",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (resp.status >= 400 && resp.status < 600) {
                                // TODO: return resp.text(). How to do?
                                throw new Error("TODO");
                            }
                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => {
                            this.error = err;
                            console.log(err);
                        });
                },
                delete: () => {
                    let params = new URLSearchParams();
                    params.append("file", this.file.filename);

                    fetch("/api/files?" + params, {
                        method: "DELETE",
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (resp.status >= 400 && resp.status < 600) {
                                // TODO: return resp.text(). How to do?
                                return Promise.reject("TODO");
                            }

                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                            return resp.json();
                        })
                        .then(resp => console.log(resp))
                        .catch(err => {
                            this.error = err;
                            console.log(err);
                        });
                }
            };
        },
        // Tags API
        tagsAPI: function() {
            return {
                // Requests
                add: (name, color) => {
                    let params = new URLSearchParams();
                    params.append("name", name);
                    params.append("color", color);

                    fetch("/api/tags", {
                        method: "POST",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            this.tagsAPI().delNewTag();
                            store.updateTags();
                            return resp.text();
                        })
                        .then(err => console.log(err)); // TODO
                },
                change: (tagID, newName, newColor) => {
                    let params = new URLSearchParams();
                    params.append("id", tagID);
                    params.append("name", newName);
                    params.append("color", newColor);

                    fetch("/api/tags", {
                        method: "PUT",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            store.updateTags();
                            return resp.text();
                        })
                        .then(err => console.log(err)); // TODO
                },
                del: tagID => {
                    let params = new URLSearchParams();
                    params.append("id", tagID);

                    fetch("/api/tags?" + params, {
                        method: "DELETE",
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            store.updateTags();
                            return resp.text();
                        })
                        .then(err => console.log(err)); // TODO
                },
                // delNewTag deletes tag from tagsNewData.newTag
                delNewTag: () => {
                    this.newTag = {};
                }
            };
        }
    }
});
