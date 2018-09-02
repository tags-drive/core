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

Vue.component("file-tag", {
    props: ["name", "color"],
    template: `
	<div :style="{ 'background-color': color }" class="tag">
		<div>{{name}}</div>
	</div>`
});

Vue.component("files", {
    props: ["file"],
    data: function() {
        return {
            hover: false
        };
    },
    template: `
	<tr
		:style="[hover ? {'background-color': 'rgba(0, 0, 0, 0.1)'} : {'background-color': 'white'} ]"
		@mouseover="hover = true;"
		@mouseleave="hover = false;"
	>
		<td v-if="file.type == 'image'" style="width: 30px;">
			<img :src="file.preview" style="width: 30px;">
		</td>
		<td v-else style="width: 30px; text-align: center;">
			<img :src="'/ext/' + file.filename.split('.').pop()" style="width: 30px;">
		</td>	
		<td style="width: 200px;">
			<div class="fileName">
				<a :href="file.origin" :title="file.filename" download>{{file.filename}}</a>
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
	</tr>`
});
