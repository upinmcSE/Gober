import { Button } from '@/components/Button';
import { Divider } from '@/components/Divider';
import { HStack } from '@/components/HStack';
import { Text } from '@/components/Text';
import { VStack } from '@/components/VStack';
import { TabBarIcon } from '@/components/navigation/TabBarIcon';
import { useAuth } from '@/context/AuthContext';
import { useOnScreenFocusCallback } from '@/hooks/useOnScreenFocusCallback';
import { eventService } from '@/services/events';
import { ticketService } from '@/services/tickets';
import { EventListData } from '@/types/event';
import { UserRole } from '@/types/user';
import { router, useNavigation } from 'expo-router';
import { useCallback, useEffect, useState } from 'react';
import { Alert, FlatList, TouchableOpacity } from 'react-native';

export default function EventsScreen() {
  const { user } = useAuth();
  const navigation = useNavigation();

  const [isLoading, setIsLoading] = useState(false);
  const [events, setEvents] = useState<EventListData>();

  function onGoToEventPage(id: number) {
    if (user?.role === UserRole.Manager) {
      router.push(`/(authed)/(tabs)/(events)/event/${id}`);
    }
  }

  async function buyTicket(id: number) {
    try {
      await ticketService.createOne(id);
      Alert.alert("Success", "Ticket purchased successfully");
      fetchEvents();
    } catch (error) {
      Alert.alert("Error", "Failed to buy ticket");
    }
  }

  const fetchEvents = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await eventService.getAll();
      setEvents(response.data[0]);
    } catch (error) {
      Alert.alert("Error", "Failed to fetch events");
    } finally {
      setIsLoading(false);
    }
  }, []);

  useOnScreenFocusCallback(fetchEvents);

  useEffect(() => {
    navigation.setOptions({
      headerTitle: "Events",
      headerRight: user?.role === UserRole.Manager ? headerRight : null,
    });
  }, [navigation, user]);

  return (
    <VStack flex={ 1 } p={ 20 } pb={ 0 } gap={ 20 }>

      <HStack alignItems='center' justifyContent='space-between'>
        <Text fontSize={ 18 } bold>{ events?.events.length } Events</Text>
      </HStack>

      <FlatList
        keyExtractor={ (item) => item.event_id.toString() }
        data={ events?.events }
        onRefresh={ fetchEvents }
        refreshing={ isLoading }
        renderItem={ ({ item: event }) => (
          <VStack
            gap={ 20 }
            p={ 20 }
            style={{
              backgroundColor: 'white',
              borderRadius: 20,
            }} 
            key={ event.event_id }>

            <TouchableOpacity onPress={ () => onGoToEventPage(event.event_id) }>
              <HStack alignItems='center' justifyContent="space-between">
                <HStack alignItems='center'>
                  <Text fontSize={ 26 } bold >{ event.title }</Text>
                  <Text fontSize={ 26 } bold > | </Text>
                  <Text fontSize={ 16 } bold >{ event.location }</Text>
                </HStack>
                { user?.role === UserRole.Manager && <TabBarIcon size={ 24 } name="chevron-forward" /> }
              </HStack>
            </TouchableOpacity>

            <Divider />

            <HStack justifyContent='space-between'>

              {/* <VStack gap={ 10 }>
                <Text bold fontSize={ 16 } color='gray'>Sold: { event.totalTicketsPurchased }</Text>
                <Text bold fontSize={ 16 } color='green'>Entered: { event.totalTicketsEntered }</Text>
              </VStack> */}

              { user?.role === UserRole.Attendee && (
                <VStack>
                  <Button
                    variant='outlined'
                    disabled={ isLoading }
                    onPress={ () => buyTicket(event.event_id) }
                  >
                    Buy Ticket
                  </Button>
                </VStack>
              ) }

            </HStack>
            <Text fontSize={ 13 } color='gray'>{ event.date }</Text>
          </VStack>

        ) }
        ItemSeparatorComponent={ () => <VStack h={ 20 } /> }
      />

    </VStack>
  );
}

const headerRight = () => {
  return (
    <TabBarIcon size={ 32 } name="add-circle-outline" onPress={ () => router.push('/(authed)/(tabs)/(events)/new') } />
  );
};