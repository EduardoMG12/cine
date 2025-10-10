module.exports = {
	"*.{js,jsx,ts,tsx}": (filenames) =>
		filenames.map((filename) => `biome format --write "${filename}"`),

	"api_v2/**/*.go": (filenames) => {
		const commands = [];
		// Format each file individually with correct path
		filenames.forEach((filename) => {
			commands.push(`gofmt -w "${filename}"`);
		});
		// Run go vet in the api_v2 directory
		commands.push("cd api_v2 && go vet ./...");
		return commands;
	},

	"flutter_app/**/*.dart": (filenames) => {
		const commands = [];
		// Format each file individually with correct path
		filenames.forEach((filename) => {
			commands.push(`flutter format "${filename}"`);
		});
		return commands;
	},
};
