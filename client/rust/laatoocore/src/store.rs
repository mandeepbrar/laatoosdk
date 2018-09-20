use action::{Action};
use reducer::{Reducer};

pub trait StoreData: Clone + Reducer {
    fn initial_state(&self) -> Self;
}

pub trait Store {
    fn dispatch(&mut self, action: Box<Action>) -> Result<(), String>;
}

#[allow(dead_code)]
pub struct GenericStore<DataType> {
    state: DataType,
    dispatching: bool,
}

impl<DataType> GenericStore<DataType> where 
    DataType: StoreData{        
    fn initial_state(&self) -> DataType {
        self.state.initial_state()
    }
    fn current_state(&self) -> DataType {
        self.state.clone()
    }
    fn dispatch(&mut self, action: Box<Action>) -> Result<(), String> {
        match self.current_state().reduce(action) {
            Ok(_) => {}
            Err(e) => {
                return Err(format!("{}", e));
            }
        }
        Ok(())
    }

}

/*impl 
{
        match self.current_state().reduce(action) {
            Ok(_) => {}
            Err(e) => {
                return Err(format!("{}", e));
            }
        }
        Ok(self.current_state().clone())
    }

    */

#[cfg(test)]
mod tests {
    struct TestStore {
        testdata: Vec<i32>,
    }
    impl Reducer for TestStore {
        type Item = TestStoreAction;
        type Error = String;
        fn reduce(&mut self, action: TestStoreAction) -> Result<Self, String> {
            Ok(self);
        }
    }

    enum TestStoreAction {
        Add(i32),
    }

    #[test]
    fn store_works() {

        assert_eq!(2 + 2, 4);
    }
}



/*
pub struct LaatooStore {

}

impl LaatooStore {

}

pub fn dispatch(&self, action: T::Action) -> Result<T::Action, String> {
        let ref dispatch = self.dispatch_chain;
        match dispatch(&self, action.clone()) {
            Err(e) => return Err(format!("Error during dispatch: {}", e)),
            _ => {}
        }

        // snapshot the active subscriptions here before calling them. This both
        // emulates the Redux.js way of doing them *and* frees up the lock so
        // that a subscription can cause another subscription; also use this
        // loop to grab the ones that are safe to remove and try to remove them
        // after this
        let (subs_to_remove, subs_to_use) = self.get_subscriptions();

        // on every subscription callback loop we gather the indexes of cancelled
        // subscriptions; if we leave a loop and have cancelled subscriptions, we'll
        // try to remove them here
        self.try_to_remove_subscriptions(subs_to_remove);

        // actually run the subscriptions here; after this method is over the subs_to_use
        // vec gets dropped, and all the Arcs of subscriptions get decremented
        for subscription in subs_to_use {
            let cb = &subscription.callback;
            cb(&self, &subscription);
        }

        Ok(action)
    }

    pub fn get_state(&self) -> T {
        self.internal_store.lock().unwrap().data.clone()
    }


    pub fn subscribe(&self, callback: SubscriptionFunc<T>) -> Arc<Subscription<T>> {
        let subscription = Arc::new(Subscription::new(callback));
        let s = subscription.clone();
        self.subscriptions.write().unwrap().push(s);
        return subscription;
    }
    */