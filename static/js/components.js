Vue.component("search-tag", {
    props: ["name", "color", "show"],
    template: `
	<div @mouseenter="show = true;" @mouseleave="show = false;" :style="{ 'background-color': color }" class="tag">
		<div>{{name}}</div>
		<i v-show="show" @click="deleteTagFromSearch(name);" class="material-icons" style="cursor: pointer; font-size: 20px;">close</i>
	</div>`,
    methods: {
        deleteTagFromSearch: function(name) {
            this.$parent.deleteTagFromSearch(name);
        }
    }
});

Vue.component("files", {
    props: ["files"],
    template: `
	<table style="width:100%;">
			<tr>
				<th></th>
				<th>Filename</th>
				<th>Tags</th>
				<th>Size (MB)</th>
				<th>Time of adding</th>
			</tr>
			<tr v-for="file in files">
				<td style="width: 30px; text-align: center;">
					<i class="material-icons">
						assignment
					</i>
				</td>
				<td style="width: 200px;">
					<div class="fileName">
						{{file.filename}}
					</div>
				</td>
				<td>
					<div style="display: flex;">
						<file-tag
							v-for="tag in file.tags"
							:name="tag.name"
							:color="tag.color">
						</file-tag>
					</div>
				</td>
				<td>{{(file.size / (1024 * 1024)).toFixed(1)}}</td>
				<td>{{file.addTime}}</td>
			</tr>
	</table>`
});

Vue.component("file-tag", {
    props: ["name", "color"],
    template: `
	<div :style="{ 'background-color': color }" class="tag">
		<div>{{name}}</div>
	</div>`
});
