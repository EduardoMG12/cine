module.exports = {
	"*.{js,jsx,ts,tsx}": (filenames) =>
		filenames.map((filename) => `biome format --write "${filename}"`),

	"api_v2/**/*.go": (filenames) => {
		const commands = [];

		// Format files
		filenames.forEach((filename) => {
			commands.push(`gofmt -w "${filename}"`);
		});

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
