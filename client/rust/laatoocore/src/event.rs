use std::fmt::Debug;
use std::any::Any;

pub trait EventProducer {
    fn register_listener(&mut self, Box<EventListener>);
}

pub trait EventListener {
    fn on_event(&self, Box<Event>);
}

pub trait Event : Debug + Any {
    fn get_type(&self) ->&str;
    fn as_any(&self) -> &dyn Any;  
}