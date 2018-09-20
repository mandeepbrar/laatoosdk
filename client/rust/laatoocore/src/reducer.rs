use std::fmt::Display;
use action::Action;
use error::Error;

pub trait Reducer {
    /// Reduce a given state based upon an action. This won't be called externally
    /// because your application will never have a reference to the state object
    /// directly. Instead, it'll be called with you call `store.dispatch`.
    fn reduce(&mut self, Box<Action>) -> Result<(), String>;
}