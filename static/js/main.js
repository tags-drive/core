// For every fetch request we have to use `credentials: "same-origin"` to pass Cookie

// store contains global data
var store = {
    state: {
        msg: "Test",
        allFiles: [],
        allTags: []
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

// Main block
var main = new Vue({
    el: "#mainBlock",
    data: {
        sharedState: store.state
    }
});

// Search bar
var searchBar = new Vue({
    el: "#searchBox",
    data: {
        sharedState: store.state,
        tagForAdding: "",
        pickedTags: [],
        text: "",
        selectedSortType: "Name",
        selectedSortOrder: "Asc",
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
            params.append("sort", this.selectedSortType.toLowerCase());
            // order
            params.append("order", this.selectedSortOrder.toLowerCase());
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

