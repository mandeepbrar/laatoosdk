use std::error;
use platform;

#[derive(Default)]
pub struct Application<'a> {
    pf: Option<Box<platform::Platform + 'a>>
}

impl <'a> Application<'a> {
    fn instance() -> Box<Application<'a>> {
        let app: Application<'static> = Application{pf: Option::None};
        Box::new(app)
    }

    fn set_platform(&mut self, pf: Box<platform::Platform>) {
        println!("set platform");
        self.pf = Option::Some(pf);
        match self.pf {
            Some(ref myval) => { println!("my func .. ");  },
            None => {println!("none .. ");}
        }
    }
}
