import { Button } from '@/components/Button';
import { Input } from '@/components/Input';
import { Text } from '@/components/Text';
import { VStack } from '@/components/VStack';
import { TabBarIcon } from '@/components/navigation/TabBarIcon';
import { useOnScreenFocusCallback } from '@/hooks/useOnScreenFocusCallback';
import { eventService } from '@/services/events';
import { Event } from '@/types/event';
import DateTimePicker from '@react-native-community/datetimepicker';
import { router, useLocalSearchParams, useNavigation } from 'expo-router';
import { useCallback, useEffect, useState } from 'react';
import { Alert } from 'react-native';

export default function EventDetailsScreen() {
  const navigation = useNavigation();
  const { id } = useLocalSearchParams();

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [eventData, setEventData] = useState<Event | null>(null);

  function updateField(field: keyof Event, value: string | Date) {
    setEventData((prev) => ({
      ...prev!,
      [field]: value,
    }));
  }

  const onDelete = useCallback(async () => {
    if (!eventData) return;
    try {
      Alert.alert("Delete Event", "Are you sure you want to delete this event?", [
        { text: "Cancel" },
        {
          text: "Delete", onPress: async () => {
            await eventService.deleteOne(Number(id));
            router.back();
          }
        },
      ]);
    } catch (error) {
      Alert.alert("Error", "Failed to delete event");
    }
  }, [eventData, id]);

  async function onSubmitChanges() {
    if (!eventData) return;
    try {
      setIsSubmitting(true);
      await eventService.updateOne(Number(id),
        eventData.title,
        eventData.location,
        eventData.date
      );
      router.back();
    } catch (error) {
      Alert.alert("Error", "Failed to update event");
    } finally {
      setIsSubmitting(false);
    }
  }

  const fetchEvent = useCallback(async () => {
    try {
      const response = await eventService.getOne(Number(id));
      setEventData(response.data[0]);
    } catch (error) {
      router.back();
    }
  }, [id, router]);

  useOnScreenFocusCallback(fetchEvent);

  useEffect(() => {
    navigation.setOptions({
      headerTitle: "",
      headerRight: () => headerRight(onDelete)
    });
  }, [navigation, onDelete]);

  return (
    <VStack m={ 20 } flex={ 1 } gap={ 30 }>

      <VStack gap={ 5 }>
        <Text ml={ 10 } fontSize={ 14 } color="gray">Name</Text>
        <Input
          value={ eventData?.title }
          onChangeText={ (value) => updateField("title", value) }
          placeholder="Title"
          placeholderTextColor="darkgray"
          h={ 48 }
          p={ 14 }
        />
      </VStack>

      <VStack gap={ 5 }>
        <Text ml={ 10 } fontSize={ 14 } color="gray">Location</Text>
        <Input
          value={ eventData?.location }
          onChangeText={ (value) => updateField("location", value) }
          placeholder="Name"
          placeholderTextColor="darkgray"
          h={ 48 }
          p={ 14 }
        />
      </VStack>

      <VStack gap={ 5 }>
        <Text ml={ 10 } fontSize={ 14 } color="gray">Date</Text>
        <DateTimePicker
          style={ {
            alignSelf: "flex-start",
          } }
          accentColor='black'
          minimumDate={ new Date() }
          value={ new Date(eventData?.date ?? new Date()) }
          mode={ "datetime" }
          onChange={ (_, date) => date && updateField("date", date) }
        />
      </VStack>

      <Button
        mt={ "auto" }
        isLoading={ isSubmitting }
        disabled={ isSubmitting }
        onPress={ onSubmitChanges }
      >
        Save Changes
      </Button>

    </VStack>
  );
}

const headerRight = (onPress: VoidFunction) => {
  return (
    <TabBarIcon size={ 30 } name="trash" onPress={ onPress } />
  );
};