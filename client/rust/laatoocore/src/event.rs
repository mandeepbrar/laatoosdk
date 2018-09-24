use std::fmt::Debug;

pub trait EventProducer {
    fn register_listener(&mut self, Box<EventListener>);
}

pub trait EventListener {
    fn on_event(&self, Box<Event>);
}

pub trait Event : Debug {
    fn get_type(&self) ->&str;
    fn get_source(&self);
    fn get_data(&self);
}