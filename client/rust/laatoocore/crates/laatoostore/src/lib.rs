pub mod store;
pub mod reducer;

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}


pub type DispatchFunc<T: Reducer> = Box<Fn(&Store<T>, T::Action) -> Result<T, String>>;

pub struct GlobalStore {

}

pub fn register_store(store: Store, Action) {

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