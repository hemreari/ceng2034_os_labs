extern crate dirs;
extern crate chrono;

use chrono::prelude::*;
use chrono::Utc;

use std::env;
use std::fs;
use std::fs::File;
use std::path::Path;
use std::os::unix::fs::MetadataExt;


fn main() {
	let path = dirs::home_dir().unwrap();
	let fullpath = path.to_str().unwrap();

	let mut os_lab_path = String::from(fullpath.to_owned());
	os_lab_path.push_str("/os_lab_0");


	//fs::create_dir(fullpath.to_owned() + "/os_lab_0").unwrap();
	match fs::create_dir(os_lab_path.clone()) {
		Err(why) => println!("Couldn't create {} : {:?}", os_lab_path, why.kind()),
		Ok(_) => println!("Folder created: {}", os_lab_path),
	};

	let root = Path::new(&os_lab_path);
	assert!(env::set_current_dir(&root).is_ok());
	println!("Changed working directory to {}", root.display());

	let mut file1_txt_path = os_lab_path.clone();
	file1_txt_path.push_str("/a.txt");

	let mut file2_txt_path = os_lab_path.clone();
	file2_txt_path.push_str("/b.txt");

	let mut file_py_path = os_lab_path.clone();
	file_py_path.push_str("/hi.py");

	let filepaths: [String; 3] = [file1_txt_path, file2_txt_path, file_py_path];

	for file_path in &filepaths {
		match File::create(file_path) {
			Err(why) => println!("Couldn't create {} : {:?}", file_path, why.kind()),
			Ok(_) => println!("File created {}", file_path),
		};
	}

	for file_path in &filepaths {
		let meta = fs::metadata(file_path).unwrap();
		let timestamp = meta.ctime();

		let naive_datetime = NaiveDateTime::from_timestamp(timestamp, 0);
		let datetime: DateTime<Utc> = DateTime::from_utc(naive_datetime, Utc);

		println!("{} last modified date: {}", file_path, datetime);

		/*
		let datetime: DateTime<Utc>::from(timestamp);
		let newdate = datetime.format("%Y-%m-%d %H:%M:%S.%f").to_string();
		let metadata = fs::metadata(file_path);
		if let Ok(time) = metadata.modified() {
			println!("{} last modifed date: {:?}",
							  file_path, time.modified());
		} else {
			println!("Couldn't print last modified date {}",
								  file_path);
		}
		let modified_date = match fs::metadata(file_path) {
			Err(why) => println!("Couldn't print last modified date {} : {:?}",
								  file_path, why.kind()),
			Ok(meta) => println!("{} last modifed date: {:?}",
							  file_path, meta.ctime()),
		};*/
	}

	//println!("file1 {}, file2 {}, filepy {}", file1_txt_path, file2_txt_path, file_py_path);

	//let mut file1_txt = match File::create(path: P)
}
